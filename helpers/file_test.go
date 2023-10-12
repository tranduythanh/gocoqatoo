package helpers

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileExists(t *testing.T) {
	Convey("Given a file path", t, func() {
		testFilePath := "test.txt" // This file should exist for the test
		defer os.Remove(testFilePath)

		Convey("When the file exists", func() {
			_, err := os.Create(testFilePath)
			So(err, ShouldBeNil)

			Convey("FileExists should return true", func() {
				So(FileExists(testFilePath), ShouldBeTrue)
			})
		})

		Convey("When the file does not exist", func() {
			So(FileExists("nonexistent.txt"), ShouldBeFalse)
		})
	})
}

func TestConvertFileToString(t *testing.T) {
	Convey("Given a file path", t, func() {
		testFilePath := "test.txt"
		defer os.Remove(testFilePath)

		Convey("When the file exists", func() {
			file, err := os.Create(testFilePath)
			So(err, ShouldBeNil)
			file.WriteString("Hello, GoConvey!")
			file.Close()

			Convey("ConvertFileToString should return file content", func() {
				content, err := ConvertFileToString(testFilePath)
				So(err, ShouldBeNil)
				So(content, ShouldEqual, "Hello, GoConvey!")
			})
		})

		Convey("When the file does not exist", func() {
			_, err := ConvertFileToString("nonexistent.txt")
			So(err, ShouldNotBeNil)
		})
	})
}
