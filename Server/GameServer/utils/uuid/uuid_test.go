package uuid

import (
	"fmt"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	//m := make(map[string]int)
	b.ResetTimer()
	l, _ := New(16)
	l2, _ := New(16)
	l3, _ := New(16)
	fmt.Println(l)
	fmt.Println(l2)
	fmt.Println(l3)
	//for i:=0;i<1000;i++ {
	//	uuid,err  := New(16)
	//	uuid,err  := New(16)
	//	//uuid,err := Generate("gate1",12)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	if _,ok := m[uuid]; ok {
	//		fmt.Println(uuid)
	//	}else{
	//		m[uuid] = 0
	//	}
	//
	//	fmt.Println(uuid)
	//}
}
