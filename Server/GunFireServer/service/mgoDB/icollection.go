package mgoDB

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICollection interface {
	CollectionName() string
	GetId() primitive.ObjectID
	SetId(id primitive.ObjectID)
}

type Filter struct {
	col    ICollection
	filter bson.D
}

func NewFilter(col ICollection) *Filter {
	return &Filter{
		col:    col,
		filter: bson.D{},
	}
}

func (o *Filter) Where(filter bson.D) *Filter {
	o.filter = append(o.filter, filter...)
	return o
}

type Finder struct {
	col     ICollection
	filter  bson.D
	records interface{}
	options []*options.FindOptions
}

func NewFinder(col ICollection) *Finder {
	return &Finder{
		col:     col,
		filter:  bson.D{},
		options: []*options.FindOptions{},
	}
}

func (o *Finder) Where(filter bson.D) *Finder {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	o.filter = append(o.filter, filter...)
	return o
}

func (o *Finder) Records(records interface{}) *Finder {
	o.records = records
	return o
}

func (o *Finder) Options(opts *options.FindOptions) *Finder {
	if o.options == nil {
		o.options = []*options.FindOptions{}
	}
	o.options = append(o.options, opts)
	return o
}

type OneFinder struct {
	col     ICollection
	record  interface{}
	filter  bson.D
	options []*options.FindOneOptions
}

func NewOneFinder(col ICollection) *OneFinder {
	return &OneFinder{
		col:     col,
		filter:  bson.D{},
		options: []*options.FindOneOptions{},
	}
}

func (o *OneFinder) Where(filter bson.D) *OneFinder {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	o.filter = append(o.filter, filter...)
	return o
}

func (o *OneFinder) Options(opts *options.FindOneOptions) *OneFinder {
	if o.options == nil {
		o.options = []*options.FindOneOptions{}
	}
	o.options = append(o.options, opts)
	return o
}
func (o *OneFinder) Record(record interface{}) *OneFinder {
	o.record = record
	return o
}
type Updater struct {
	col     ICollection
	filter  bson.D
	update  bson.D
	push    bson.D
	inc    bson.D
	options []*options.UpdateOptions
}

func NewUpdater(col ICollection) *Updater {
	return &Updater{
		col:     col,
		filter:  bson.D{},
		update:  bson.D{},
		options: []*options.UpdateOptions{},
		inc:  bson.D{},
	}
}

func (o *Updater) Where(filter bson.D) *Updater {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	o.filter = append(o.filter, filter...)
	return o
}

func (o *Updater) Update(update bson.D) *Updater {
	o.update = append(o.update, update...)
	return o
}

func (o *Updater) Push(push bson.D) *Updater {
	o.push = append(o.push, push...)
	return o
}
func (o *Updater) Inc(push bson.D) *Updater {
	o.inc = append(o.inc, push...)
	return o
}

func (o *Updater) Options(opts *options.UpdateOptions) *Updater {
	if o.options == nil {
		o.options = []*options.UpdateOptions{}
	}
	o.options = append(o.options, opts)
	return o
}

type Deleter struct {
	col     ICollection
	filter  bson.D
	options []*options.DeleteOptions
}

func NewDeleter(col ICollection) *Deleter {
	return &Deleter{
		col:     col,
		filter:  bson.D{},
		options: []*options.DeleteOptions{},
	}
}

func (o *Deleter) Where(filter bson.D) *Deleter {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	o.filter = append(o.filter, filter...)
	return o
}

func (o *Deleter) Options(opts *options.DeleteOptions) *Deleter {
	if o.options == nil {
		o.options = []*options.DeleteOptions{}
	}
	o.options = append(o.options, opts)
	return o
}

type Aggregator struct {
	col      ICollection
	pipeline bson.A
	options  []*options.AggregateOptions
	records  interface{}
}

func NewAggregator(col ICollection) *Aggregator {
	return &Aggregator{
		col:      col,
		pipeline: bson.A{},
		options:  []*options.AggregateOptions{},
	}
}

func (o *Aggregator) Stage(stage bson.D) *Aggregator {
	if o.pipeline == nil {
		o.pipeline = bson.A{}
	}
	o.pipeline = append(o.pipeline, stage)
	return o
}

func (o *Aggregator) Options(opts *options.AggregateOptions) *Aggregator {
	if o.options == nil {
		o.options = []*options.AggregateOptions{}
	}
	o.options = append(o.options, opts)
	return o
}

func (o *Aggregator) Records(records interface{}) *Aggregator {
	o.records = records
	return o
}

type Counter struct {
	col     ICollection
	filter  bson.D
	options []*options.CountOptions
}

func NewCounter(col ICollection) *Counter {
	return &Counter{
		col:    col,
		filter: bson.D{},
	}
}

func (o *Counter) Where(filter bson.D) *Counter {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	o.filter = append(o.filter, filter...)
	return o
}

func (o *Counter) Options(opts *options.CountOptions) *Counter {
	if o.options == nil {
		o.options = []*options.CountOptions{}
	}
	o.options = append(o.options, opts)
	return o
}

type EstimateCounter struct {
	col     ICollection
	options []*options.EstimatedDocumentCountOptions
}

func NewEstimateCounter(col ICollection) *EstimateCounter {
	return &EstimateCounter{
		col: col,
	}
}

func (o *EstimateCounter) Options(opts *options.EstimatedDocumentCountOptions) *EstimateCounter {
	if o.options == nil {
		o.options = []*options.EstimatedDocumentCountOptions{}
	}
	o.options = append(o.options, opts)
	return o
}
