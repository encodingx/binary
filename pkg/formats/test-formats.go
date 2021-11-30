package formats

type TestFormat0 struct {
	TestWord0 `word:"24"`
	TestWord1 `word:"32"`
	TestWord2 `word:"40"`
}

type TestFormat1Bad struct {
	TestWord3Bad `word:"24"`
	TestWord1    `word:"32"`
	TestWord2    `word:"40"`
}

type TestFormat2Bad struct {
	TestWord0    `word:"24"`
	TestWord4Bad `word:"32"`
	TestWord2    `word:"40"`
}

type TestFormat3Bad struct {
	TestWord0    `word:"24"`
	TestWord1    `word:"32"`
	TestWord5Bad `word:"40"`
}

type TestFormat4Bad struct {
	TestWord6Bad `word:"24"`
	TestWord1    `word:"32"`
	TestWord2    `word:"40"`
}

type TestFormat5Bad struct {
	TestWord0    `word:"24"`
	TestWord7Bad `word:"32"`
	TestWord2    `word:"40"`
}

type TestFormat6Bad struct {
	TestWord0    `word:"24"`
	TestWord1    `word:"32"`
	TestWord8Bad `word:"40"`
}

type TestFormat7Bad struct {
	TestWord9Bad `word:"20"` // word length not a multiple of 8
	TestWord1    `word:"32"`
	TestWord2    `word:"40"`
}

type TestFormat8Bad struct {
	TestWord0 `word:"24"`
	TestWord1 // missing struct tag
	TestWord2 `word:"40"`
}

type TestFormat9Bad struct {
	TestWord0 `word:"24"`
	TestWord1 `word:"32"`
	TestWord2 `worm:"40"` // misspelt struct tag
}

type TestFormat10Bad struct {
	TestWord0 `word:"16"` // word length different from actual
	TestWord1 `word:"32"`
	TestWord2 `word:"40"`
}

type TestFormat11Bad struct {
	TestWord0 `word:"24"`
	TestWord1 `word:"72"` // word size exceeds limit
	TestWord2 `word:"40"`
}

type TestFormat12Bad struct {
}

type TestFormat13Bad struct {
	TestWord0  `word:"24"`
	TestWord1  `word:"32"`
	TestField2 byte `bitfield:"8,0"` // bitfield in format
}

type TestWord0 struct { //               2   1         0
	//                                   321098765432109876543210
	TestField0 byte `bitfield:"8,16"` // |------|
	TestField1 byte `bitfield:"8,8"`  //         |------|
	TestField2 byte `bitfield:"8,0"`  //                 |------|
}

type TestWord1 struct { //               3 2         1         0
	//                                   10987654321098765432109876543210
	TestField0 byte `bitfield:"5,27"` // |---|
	TestField1 byte `bitfield:"7,20"` //      |-----|
	TestField2 byte `bitfield:"9,11"` //             |-------|
	TestField3 byte `bitfield:"11,0"` //                      |---------|
}

type TestWord2 struct { //                3         2         1         0
	//                                    9876543210987654321098765432109876543210
	TestField0 byte `bitfield:"3,37"`  // |-|
	TestField1 byte `bitfield:"10,27"` //    |--------|
	TestField2 byte `bitfield:"6,21"`  //              |----|
	TestField3 byte `bitfield:"13,8"`  //                    |-----------|
	TestField4 byte `bitfield:"8,0"`   //                                 |------|
}

type TestWord3Bad struct { //            2   1         0
	// gap between fields                321098765432109876543210
	TestField0 byte `bitfield:"6,18"` // |----|
	TestField1 byte `bitfield:"8,8"`  //         |------|
	TestField2 byte `bitfield:"8,0"`  //                 |------|
}

type TestWord4Bad struct { //            3 2         1         0
	// overlapping fields                10987654321098765432109876543210
	TestField0 byte `bitfield:"5,27"` // |---|
	TestField1 byte `bitfield:"9,18"` //      |-------|
	TestField2 byte `bitfield:"9,11"` //             |-------|
	TestField3 byte `bitfield:"11,0"` //                      |---------|
}

type TestWord5Bad struct { //             3         2         1         0
	// gap and overlap between fields     9876543210987654321098765432109876543210
	TestField0 byte `bitfield:"3,29"`  // |-|
	TestField1 byte `bitfield:"10,19"` //    |--------|
	TestField2 byte `bitfield:"6,11"`  //                |----|
	TestField3 byte `bitfield:"13,0"`  //                    |-----------|
	TestField4 byte `bitfield:"8,0"`   //                                 |------|
}

type TestWord6Bad struct { //           2   1         0
	// missing field                    321098765432109876543210
	TestField1 byte `bitfield:"8,8"` //         |------|
	TestField2 byte `bitfield:"8,0"` //                 |------|
}

type TestWord7Bad struct { //            3 2         1         0
	// missing struct tag                10987654321098765432109876543210
	TestField0 byte `bitfield:"5,27"` // |---|
	TestField1 byte
	TestField2 byte `bitfield:"9,11"` //             |-------|
	TestField3 byte `bitfield:"11,0"` //                      |---------|
}

type TestWord8Bad struct { //             3         2         1         0
	// repeated fields                    9876543210987654321098765432109876543210
	TestField0 byte `bitfield:"3,37"`  // |-|
	TestField1 byte `bitfield:"10,27"` //    |--------|
	TestField2 byte `bitfield:"6,21"`  //              |----|
	TestField3 byte `bitfield:"6,21"`  //              |----|
	TestField4 byte `bitfield:"13,8"`  //                    |-----------|
	TestField5 byte `bitfield:"8,0"`   //                                 |------|
}

type TestWord9Bad struct { //            1         0
	// total length not mulitple of 8    98765432109876543210
	TestField0 byte `bitfield:"8,12"` // |------|
	TestField1 byte `bitfield:"8,4"`  //         |------|
	TestField2 byte `bitfield:"4,0"`  //                 |--|
}

type TestWord10Bad struct {
}

type TestWord11Bad struct {
	TestField0 string `bitfield:"8,16"` // unsupported field type
	TestField1 byte   `bitfield:"8,8"`
	TestField2 byte   `bitfield:"8,0"`
}

type TestWord12Bad struct {
	TestField0 byte `bitfield:"5,27"`
	TestField1 byte `bitfield:"7,20"`
	TestField2 byte `bitfield:"9,11"` // field size overflows type
	TestField3 uint `bitfield:"11,0"`
}
