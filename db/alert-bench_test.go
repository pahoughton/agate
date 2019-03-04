/* 2019-02-21 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package db

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func BenchmarkAlertAdd(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		today := time.Now()
		data := make(map[string]string)
		for i := 0; i < b.N; i += 1 {
			data[fmt.Sprintf("k-%03d",i)] = strconv.Itoa(i)
		}
		b.ResetTimer()
		for k, v := range data {
			db.AlertAdd(today,[]byte(k),[]byte(v))
		}
	})
}
func BenchmarkAlertAddGet(b *testing.B) {
	benchit(b,func (b *testing.B,db *DB) {
		today := time.Now()
		data := make(map[string]string,b.N)
		for i := 0; i < b.N; i += 1 {
			data[fmt.Sprintf("k-%09d",i)] = strconv.Itoa(i)
		}
		b.ResetTimer()
		for k, v := range data {
			db.AlertAdd(today,[]byte(k),[]byte(v))
		}
		for k, _ := range data {
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
