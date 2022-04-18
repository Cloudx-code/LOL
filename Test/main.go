package main

import "fmt"

func main() {
	a, b, c := 43, 53, 59
	d := ((float64(a) + float64(c)) / float64(b))
	fmt.Println(d)

	//var s []int = make( []int,0 )
	//for i :=0;i<5;i++{
	//	s = append( s,i )
	//}
	//for _,i := range s {
	//	fmt.Println( i )
	//}
}
