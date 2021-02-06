package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"testing"
)

// TestLineSplitterWrite ensures various length data is correctly broken up
func TestLineSplitterWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// Parameters.
		p []byte
		// Expected results.
		want string
	}{
		{
			"test_archive.zip",
			[]byte{0x50, 0x4b, 0x03, 0x04, 0x14, 0x00, 0x09, 0x00, 0x08, 0x00, 0x66, 0x6d, 0xcf, 0x48, 0xb4, 0xf8,
				0x71, 0xdd, 0x53, 0x01, 0x00, 0x00, 0xd0, 0x01, 0x00, 0x00, 0x2d, 0x00, 0x1c, 0x00, 0x64, 0x65,
				0x61, 0x6c, 0x5f, 0x36, 0x39, 0x30, 0x30, 0x32, 0x34, 0x5f, 0x32, 0x30, 0x31, 0x36, 0x2d, 0x30,
				0x36, 0x2d, 0x31, 0x35, 0x2d, 0x31, 0x33, 0x34, 0x33, 0x5f, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74,
				0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x63, 0x73, 0x76, 0x55, 0x54, 0x09, 0x00, 0x03,
				0x5f, 0x4d, 0x61, 0x57, 0x5e, 0x4d, 0x61, 0x57, 0x75, 0x78, 0x0b, 0x00, 0x01, 0x04, 0xe8, 0x03,
				0x00, 0x00, 0x04, 0xe8, 0x03, 0x00, 0x00, 0xce, 0xcf, 0x3e, 0x09, 0x39, 0x2c, 0x41, 0x4e, 0x46,
				0x01, 0xc6, 0x5c, 0x96, 0xd7, 0x7a, 0x5a, 0x3c, 0xf7, 0xa7, 0x4a, 0xfa, 0x62, 0xed, 0xc8, 0x34,
				0x14, 0xff, 0x69, 0x13, 0x1c, 0xb6, 0xe2, 0x97, 0x94, 0xf9, 0xbe, 0x4b, 0x37, 0x52, 0x62, 0x45,
				0xe0, 0xbf, 0xda, 0xd7, 0x5a, 0xe9, 0xee, 0xe1, 0x2f, 0x33, 0x9e, 0x1e, 0xc9, 0x99, 0x48, 0x80,
				0x4b, 0x15, 0xf5, 0x61, 0x5a, 0x21, 0x66, 0xce, 0x5f, 0x1c, 0x6a, 0xbb, 0x91, 0x65, 0xa3, 0x0f,
				0xeb, 0x5f, 0xc4, 0xa8, 0x9f, 0x82, 0x11, 0x1d, 0xf2, 0x9c, 0x9d, 0x94, 0x1f, 0xbd, 0x80, 0x1c,
				0x8a, 0xa5, 0x80, 0xae, 0x3f, 0x40, 0x50, 0x88, 0x4b, 0x5b, 0x67, 0x75, 0xc6, 0x9e, 0x6d, 0x23,
				0x84, 0xe2, 0xa2, 0x79, 0x69, 0x61, 0xc4, 0x03, 0x1a, 0xc4, 0xc4, 0x4b, 0xf9, 0xbe, 0xe1, 0x5e,
				0xe1, 0xd8, 0xb0, 0xf5, 0x1e, 0xc8, 0xd6, 0xb0, 0x34, 0x22, 0x87, 0xac, 0x65, 0xa7, 0x0a, 0x73,
				0x72, 0x5e, 0x33, 0x75, 0x81, 0xef, 0x7e, 0x91, 0xdf, 0x04, 0x17, 0x90, 0xeb, 0xa9, 0xfa, 0x92,
				0x5e, 0xb4, 0x0f, 0xe3, 0x0a, 0x68, 0x83, 0xc7, 0xc7, 0x75, 0x3a, 0xb4, 0xb7, 0x81, 0x17, 0x5f,
				0x27, 0xae, 0xe0, 0x5b, 0x0a, 0x90, 0xad, 0x6e, 0x8c, 0xc6, 0x01, 0x0a, 0xb7, 0xe7, 0xba, 0xfd,
				0x1d, 0x7f, 0x07, 0xe1, 0xd8, 0xe6, 0x61, 0x33, 0x22, 0x63, 0xe3, 0x70, 0x49, 0xd3, 0x70, 0x4d,
				0x11, 0x06, 0x05, 0x32, 0xd9, 0x5e, 0xfd, 0x72, 0x64, 0xef, 0x5e, 0xf8, 0xd4, 0x98, 0xa6, 0xe8,
				0xe1, 0x6f, 0x87, 0xd5, 0x05, 0x96, 0xf3, 0x3a, 0x60, 0x87, 0x7e, 0x94, 0x69, 0xcd, 0x69, 0x7f,
				0x8b, 0x8e, 0xbb, 0x2d, 0xeb, 0xa1, 0x86, 0x2f, 0xe9, 0x6d, 0x87, 0x36, 0x2e, 0xe4, 0xe0, 0xec,
				0x68, 0x70, 0x8e, 0x7e, 0x26, 0xd7, 0x73, 0xf7, 0x07, 0xb9, 0x5c, 0xa0, 0x08, 0x51, 0xc9, 0x50,
				0x7c, 0xb0, 0xef, 0xad, 0x8a, 0x0d, 0x3d, 0x5d, 0x6a, 0x7d, 0x6c, 0x59, 0x36, 0x53, 0x04, 0xaa,
				0x5b, 0x2e, 0x63, 0x5b, 0xd5, 0x00, 0x06, 0x84, 0xbc, 0x6c, 0x3e, 0xf5, 0xc7, 0x52, 0x1d, 0x48,
				0xc5, 0x61, 0x1a, 0x69, 0x2f, 0xba, 0x83, 0x34, 0xe1, 0xda, 0xb3, 0x3f, 0xd6, 0x31, 0x89, 0xd2,
				0x10, 0xb2, 0xba, 0x7e, 0x9d, 0xab, 0x4b, 0xf7, 0x33, 0xc5, 0x06, 0x5e, 0x91, 0x5f, 0xa6, 0xfb,
				0xe3, 0x38, 0xbd, 0x39, 0x80, 0xf5, 0xb2, 0x4b, 0x6d, 0xdb, 0x50, 0x4b, 0x07, 0x08, 0xb4, 0xf8,
				0x71, 0xdd, 0x53, 0x01, 0x00, 0x00, 0xd0, 0x01, 0x00, 0x00, 0x50, 0x4b, 0x01, 0x02, 0x1e, 0x03,
				0x14, 0x00, 0x09, 0x00, 0x08, 0x00, 0x66, 0x6d, 0xcf, 0x48, 0xb4, 0xf8, 0x71, 0xdd, 0x53, 0x01,
				0x00, 0x00, 0xd0, 0x01, 0x00, 0x00, 0x2d, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00,
				0x00, 0x00, 0xb4, 0x81, 0x00, 0x00, 0x00, 0x00, 0x64, 0x65, 0x61, 0x6c, 0x5f, 0x36, 0x39, 0x30,
				0x30, 0x32, 0x34, 0x5f, 0x32, 0x30, 0x31, 0x36, 0x2d, 0x30, 0x36, 0x2d, 0x31, 0x35, 0x2d, 0x31,
				0x33, 0x34, 0x33, 0x5f, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72,
				0x73, 0x2e, 0x63, 0x73, 0x76, 0x55, 0x54, 0x05, 0x00, 0x03, 0x5f, 0x4d, 0x61, 0x57, 0x75, 0x78,
				0x0b, 0x00, 0x01, 0x04, 0xe8, 0x03, 0x00, 0x00, 0x04, 0xe8, 0x03, 0x00, 0x00, 0x50, 0x4b, 0x05,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x73, 0x00, 0x00, 0x00, 0xca, 0x01, 0x00,
				0x00, 0x00, 0x00,
			},
			"UEsDBBQACQAIAGZtz0i0+HHdUwEAANABAAAtABwAZGVhbF82OTAwMjRfMjAx\r\n" +
				"Ni0wNi0xNS0xMzQzX2xhdGVzdF9vcmRlcnMuY3N2VVQJAANfTWFXXk1hV3V4\r\n" +
				"CwABBOgDAAAE6AMAAM7PPgk5LEFORgHGXJbXelo896dK+mLtyDQU/2kTHLbi\r\n" +
				"l5T5vks3UmJF4L/a11rp7uEvM54eyZlIgEsV9WFaIWbOXxxqu5Flow/rX8So\r\n" +
				"n4IRHfKcnZQfvYAciqWArj9AUIhLW2d1xp5tI4TionlpYcQDGsTES/m+4V7h\r\n" +
				"2LD1HsjWsDQih6xlpwpzcl4zdYHvfpHfBBeQ66n6kl60D+MKaIPHx3U6tLeB\r\n" +
				"F18nruBbCpCtbozGAQq357r9HX8H4djmYTMiY+NwSdNwTREGBTLZXv1yZO9e\r\n" +
				"+NSYpujhb4fVBZbzOmCHfpRpzWl/i467Leuhhi/pbYc2LuTg7Ghwjn4m13P3\r\n" +
				"B7lcoAhRyVB8sO+tig09XWp9bFk2UwSqWy5jW9UABoS8bD71x1IdSMVhGmkv\r\n" +
				"uoM04dqzP9YxidIQsrp+natL9zPFBl6RX6b74zi9OYD1sktt21BLBwi0+HHd\r\n" +
				"UwEAANABAABQSwECHgMUAAkACABmbc9ItPhx3VMBAADQAQAALQAYAAAAAAAB\r\n" +
				"AAAAtIEAAAAAZGVhbF82OTAwMjRfMjAxNi0wNi0xNS0xMzQzX2xhdGVzdF9v\r\n" +
				"cmRlcnMuY3N2VVQFAANfTWFXdXgLAAEE6AMAAAToAwAAUEsFBgAAAAABAAEA\r\n" +
				"cwAAAMoBAAAAAA==",
		},
		{
			"test",
			[]byte("test"),
			"dGVzdA==",
		},
		{
			"sentance",
			[]byte("A man may fight for many things. His country, his principles, his friends. " +
				"The glistening tear on the cheek of a golden child. But personally, I'd mud-wrestle " +
				"my own mother for a ton of cash, an amusing clock and a sack of French porn."),
			"QSBtYW4gbWF5IGZpZ2h0IGZvciBtYW55IHRoaW5ncy4gSGlzIGNvdW50cnks\r\n" +
				"IGhpcyBwcmluY2lwbGVzLCBoaXMgZnJpZW5kcy4gVGhlIGdsaXN0ZW5pbmcg\r\n" +
				"dGVhciBvbiB0aGUgY2hlZWsgb2YgYSBnb2xkZW4gY2hpbGQuIEJ1dCBwZXJz\r\n" +
				"b25hbGx5LCBJJ2QgbXVkLXdyZXN0bGUgbXkgb3duIG1vdGhlciBmb3IgYSB0\r\n" +
				"b24gb2YgY2FzaCwgYW4gYW11c2luZyBjbG9jayBhbmQgYSBzYWNrIG9mIEZy\r\n" +
				"ZW5jaCBwb3JuLg==",
		},
		{
			"test.png",
			[]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0x01, 0x2c, 0x00, 0x00, 0x01, 0x2c, 0x08, 0x02, 0x00, 0x00, 0x00, 0xf6, 0x1f, 0x19,
				0x22, 0x00, 0x00, 0x03, 0x97, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0xed, 0xd9, 0x31, 0x8a, 0xc3,
				0x40, 0x14, 0x44, 0xc1, 0x1e, 0xe3, 0xfb, 0x5f, 0x59, 0x8a, 0x9d, 0x09, 0x1c, 0xbc, 0x40, 0x55,
				0x6c, 0xb4, 0x20, 0x70, 0xf2, 0x68, 0x98, 0x7f, 0xb6, 0x6b, 0xbb, 0xce, 0xef, 0xdf, 0xb6, 0xf3,
				0xe8, 0x9f, 0xf3, 0xad, 0x6f, 0x7d, 0xfb, 0xe7, 0xb7, 0x9f, 0x01, 0xa9, 0xef, 0x4e, 0xfd, 0x13,
				0xe0, 0xdd, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0xbc, 0x8e, 0x42,
				0xcc, 0x12, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84,
				0x10, 0x73, 0xa2, 0x80, 0x98, 0x25, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62,
				0x22, 0x84, 0x98, 0x08, 0x21, 0xe6, 0x44, 0x01, 0x31, 0x4b, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84,
				0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0xaf, 0xa3, 0x10, 0xb3, 0x84, 0x10, 0x13, 0x21, 0xc4,
				0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x9c, 0x28, 0x20, 0x66, 0x09,
				0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x39,
				0x51, 0x40, 0xcc, 0x12, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42,
				0xcc, 0xeb, 0x28, 0xc4, 0x2c, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13,
				0x21, 0xc4, 0x44, 0x08, 0x31, 0x27, 0x0a, 0x88, 0x59, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84,
				0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x4e, 0x14, 0x10, 0xb3, 0x84, 0x10, 0x13,
				0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0xf3, 0x3a, 0x0a, 0x31, 0x4b, 0x08,
				0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0xcc, 0x89,
				0x02, 0x62, 0x96, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62,
				0x22, 0x84, 0x98, 0x13, 0x05, 0xc4, 0x2c, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84,
				0x10, 0x13, 0x21, 0xc4, 0xbc, 0x8e, 0x42, 0xcc, 0x12, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4,
				0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x73, 0xa2, 0x80, 0x98, 0x25, 0x84, 0x98, 0x08,
				0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0xe6, 0x44, 0x01, 0x31,
				0x4b, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0xaf, 0xa3,
				0x10, 0xb3, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13,
				0x21, 0xc4, 0x9c, 0x28, 0x20, 0x66, 0x09, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84,
				0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x39, 0x51, 0x40, 0xcc, 0x12, 0x42, 0x4c, 0x84, 0x10, 0x13,
				0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0xcc, 0xeb, 0x28, 0xc4, 0x2c, 0x21, 0xc4, 0x44, 0x08,
				0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x27, 0x0a, 0x88, 0x59,
				0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62,
				0x4e, 0x14, 0x10, 0xb3, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84,
				0x10, 0xf3, 0x3a, 0x0a, 0x31, 0x4b, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4,
				0x44, 0x08, 0x31, 0x11, 0x42, 0xcc, 0x89, 0x02, 0x62, 0x96, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08,
				0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x13, 0x05, 0xc4, 0x2c, 0x21, 0xc4,
				0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0xbc, 0x8e, 0x42, 0xcc, 0x12,
				0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x73,
				0xa2, 0x80, 0x98, 0x25, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84,
				0x98, 0x08, 0x21, 0xe6, 0x44, 0x01, 0x31, 0x4b, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13,
				0x21, 0xc4, 0x44, 0x08, 0x31, 0xaf, 0xa3, 0x10, 0xb3, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08,
				0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x9c, 0x28, 0x20, 0x66, 0x09, 0x21, 0x26,
				0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x39, 0x51, 0x40,
				0xcc, 0x12, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0xcc, 0xeb,
				0x28, 0xc4, 0x2c, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4,
				0x44, 0x08, 0x31, 0x27, 0x0a, 0x88, 0x59, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08,
				0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x4e, 0x14, 0x10, 0xb3, 0x84, 0x10, 0x13, 0x21, 0xc4,
				0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0xf3, 0x3a, 0x0a, 0x31, 0x4b, 0x08, 0x31, 0x11,
				0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0xcc, 0x89, 0x02, 0x62,
				0x96, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84,
				0x98, 0x13, 0x05, 0xc4, 0x2c, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13,
				0x21, 0xc4, 0xbc, 0x8e, 0x42, 0xcc, 0x12, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08,
				0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x73, 0xa2, 0x80, 0x98, 0x25, 0x84, 0x98, 0x08, 0x21, 0x26,
				0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0xe6, 0x44, 0x01, 0x31, 0x4b, 0x08,
				0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0xaf, 0xa3, 0x10, 0xb3,
				0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4,
				0x9c, 0x28, 0x20, 0x66, 0x09, 0x21, 0x26, 0x42, 0x88, 0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08,
				0x21, 0x26, 0x42, 0x88, 0x39, 0x51, 0x40, 0xcc, 0x12, 0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4,
				0x44, 0x08, 0x31, 0x11, 0x42, 0xcc, 0xeb, 0x28, 0xc4, 0x2c, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x11,
				0x42, 0x4c, 0x84, 0x10, 0x13, 0x21, 0xc4, 0x44, 0x08, 0x31, 0x27, 0x0a, 0x88, 0x59, 0x42, 0x88,
				0x89, 0x10, 0x62, 0x22, 0x84, 0x98, 0x08, 0x21, 0x26, 0x42, 0x88, 0xdd, 0x07, 0xb4, 0x05, 0x5f,
				0x21, 0xcb, 0x54, 0xd1, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
			},
			"iVBORw0KGgoAAAANSUhEUgAAASwAAAEsCAIAAAD2HxkiAAADl0lEQVR4nO3Z\r\n" +
				"MYrDQBREwR7j+19Zip0JHLxAVWy0IHDyaJh/tmu7zu/ftvPon/Otb33757ef\r\n" +
				"AanvTv0T4N1ECDERQkyEEBMhxLyOQswSQkyEEBMhxEQIMRFCTIQQc6KAmCWE\r\n" +
				"mAghJkKIiRBiIoSYCCHmRAExSwgxEUJMhBATIcRECDGvoxCzhBATIcRECDER\r\n" +
				"QkyEEBMhxJwoIGYJISZCiIkQYiKEmAghJkKIOVFAzBJCTIQQEyHERAgxEULM\r\n" +
				"6yjELCHERAgxEUJMhBATIcRECDEnCohZQoiJEGIihJgIISZCiIkQYk4UELOE\r\n" +
				"EBMhxEQIMRFCTIQQ8zoKMUsIMRFCTIQQEyHERAgxEULMiQJilhBiIoSYCCEm\r\n" +
				"QoiJEGIihJgTBcQsIcRECDERQkyEEBMhxLyOQswSQkyEEBMhxEQIMRFCTIQQ\r\n" +
				"c6KAmCWEmAghJkKIiRBiIoSYCCHmRAExSwgxEUJMhBATIcRECDGvoxCzhBAT\r\n" +
				"IcRECDERQkyEEBMhxJwoIGYJISZCiIkQYiKEmAghJkKIOVFAzBJCTIQQEyHE\r\n" +
				"RAgxEULM6yjELCHERAgxEUJMhBATIcRECDEnCohZQoiJEGIihJgIISZCiIkQ\r\n" +
				"Yk4UELOEEBMhxEQIMRFCTIQQ8zoKMUsIMRFCTIQQEyHERAgxEULMiQJilhBi\r\n" +
				"IoSYCCEmQoiJEGIihJgTBcQsIcRECDERQkyEEBMhxLyOQswSQkyEEBMhxEQI\r\n" +
				"MRFCTIQQc6KAmCWEmAghJkKIiRBiIoSYCCHmRAExSwgxEUJMhBATIcRECDGv\r\n" +
				"oxCzhBATIcRECDERQkyEEBMhxJwoIGYJISZCiIkQYiKEmAghJkKIOVFAzBJC\r\n" +
				"TIQQEyHERAgxEULM6yjELCHERAgxEUJMhBATIcRECDEnCohZQoiJEGIihJgI\r\n" +
				"ISZCiIkQYk4UELOEEBMhxEQIMRFCTIQQ8zoKMUsIMRFCTIQQEyHERAgxEULM\r\n" +
				"iQJilhBiIoSYCCEmQoiJEGIihJgTBcQsIcRECDERQkyEEBMhxLyOQswSQkyE\r\n" +
				"EBMhxEQIMRFCTIQQc6KAmCWEmAghJkKIiRBiIoSYCCHmRAExSwgxEUJMhBAT\r\n" +
				"IcRECDGvoxCzhBATIcRECDERQkyEEBMhxJwoIGYJISZCiIkQYiKEmAghJkKI\r\n" +
				"OVFAzBJCTIQQEyHERAgxEULM6yjELCHERAgxEUJMhBATIcRECDEnCohZQoiJ\r\n" +
				"EGIihJgIISZCiN0HtAVfIctU0QAAAABJRU5ErkJggg==",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			w := &lineSplitter{w: &buf, maxLen: maxLineLen}

			encoder := base64.NewEncoder(base64.StdEncoding, w)
			_, err := encoder.Write(tt.p)
			encoder.Close()

			if err != nil {
				t.Fatalf("%q. base64LineWriter.Write() error = %v", tt.name, err)
			}

			if buf.String() != tt.want {
				t.Errorf("%q. base64LineWriter.Write() = \n%v\n, want \n%v\n", tt.name, buf.String(), tt.want)
			}
		})
	}
}

func TestChunkedWrites(t *testing.T) {
	s := "a 21 character string"

	var buf bytes.Buffer
	w := &lineSplitter{w: &buf, maxLen: maxLineLen}

	for i := 0; i < 20; i++ {
		if n, err := w.Write([]byte(s)); n != len(s) || err != nil {
			t.Fatalf("wrote %d, err %v", n, err)
		}
	}

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		if len(scanner.Text()) > maxLineLen {
			t.Errorf("got linelength = %d want <= %d\n", len(scanner.Text()), maxLineLen)
		}
	}
}