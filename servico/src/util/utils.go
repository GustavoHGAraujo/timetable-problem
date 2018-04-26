package util

import l "log"
import "fmt"

const tag_error = `ERROR`
const tag_debug = `DEBUG`

func log(level, tag string, values ...interface {}) {
  l.Printf("%s\t%s\t%s\n", level, tag, fmt.Sprint(values))
}

func logf(level, tag, format string, values ...interface {}) {
  l.Printf("1:%s\t2:%s\t3:%+v\n", level, tag, fmt.Sprintf(format, values...))
}

func LogD(tag string, values ...interface {}) {
  log(tag_debug, tag, values)
}

func LogDf(tag, format string, values ...interface {}) {
  logf(tag_debug, tag, format, values)
}

func LogE(tag string, values ...interface {}) {
  log(tag_error, tag, values)
}

func LogEf(tag, format string, values ...interface {}) {
  logf(tag_error, tag, format, values)
}
