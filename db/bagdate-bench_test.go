/* 2019-02-21 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"fmt"
	"strconv"
	"testing"
)


func BenchmarkBagDateAdd(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		data := make(map[string]string)
		for i := 0; i < b.N; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		b.ResetTimer()
		for k, v := range data {
			db.BagDateAdd(tnow,tbag,[]byte(k),[]byte(v))
		}
	})
}
func BenchmarkBagDateAddGet(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		data := make(map[string]string,b.N)
		for i := 0; i < b.N; i += 1 {
			data[fmt.Sprintf("k-%09d",i)] = strconv.Itoa(i)
		}
		b.ResetTimer()
		for k, v := range data {
			db.BagDateAdd(tnow,tbag,[]byte(k),[]byte(v))
		}
		for k, _ := range data {
			db.BagDateGet(tnow,tbag,[]byte(k))
		}

	})
}
func BenchmarkBagDateDel(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		data := make(map[string]string)
		for i := 0; i < b.N; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		for k, v := range data {
			db.BagDateAdd(tnow,tbag,[]byte(k),[]byte(v))
		}
		b.ResetTimer()
		for k, _ := range data {
			db.BagDateDel(tnow,tbag,[]byte(k))
		}

	})
}
