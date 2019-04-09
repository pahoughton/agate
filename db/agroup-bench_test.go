/* 2019-02-21 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

/*
import (
	"fmt"
	"testing"
)

func BenchmarkAGroupAdd(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		data := make(map[int]string,b.N)
		for i := 0; i < b.N; i += 1 {
			data[i] = fmt.Sprintf(`{"json": "k-%03d"}`,i)
		}
		b.ResetTimer()
		for k, v := range data {
			db.AGroupAdd([]byte(v), k % 3 == 0)
		}
	})
}
*/

/*
func BenchmarkAlertGet(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		today := time.Now()
		print(fmt.Sprintf("brun %d %v\n",b.N,today))
		data := make(map[string]string,b.N)
		for i := 0; i < b.N; i += 1 {
			data[fmt.Sprintf("k-%09d",i)] = strconv.Itoa(i)
		}
		print(fmt.Sprintf("brun add %d %v\n",b.N,time.Now()))
		for k, v := range data {
			// print(" add " + k + "\n")
			db.AlertAdd(today,[]byte(k),[]byte(v))
		}
		b.ResetTimer()
		print(fmt.Sprintf("brun get %d %v\n",b.N,time.Now()))
		for k, _ := range data {
			// print(" get " + k + "\n")
			db.AlertGet(today,[]byte(k))
		}

	})
}
func BenchmarkAlertDel(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		today := time.Now()
		data := make(map[string]string)
		for i := 0; i < b.N; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		for k, v := range data {
			db.AlertAdd(today,[]byte(k),[]byte(v))
		}
		b.ResetTimer()
		for k, _ := range data {
			db.AlertDel(today,[]byte(k))
		}

	})
}
*/
