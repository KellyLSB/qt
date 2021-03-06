package templater

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/therecipe/qt/internal/binding/parser"
)

func HTemplate(m string) []byte {
	var bb = new(bytes.Buffer)
	defer bb.Reset()

	if m != parser.MOC {
		m = "Qt" + m
	}

	//header
	{
		fmt.Fprintf(bb, "%v\n\n", buildTags(m))

		fmt.Fprint(bb, "#pragma once\n\n")

		fmt.Fprintf(bb, "#ifndef GO_%v_H\n", strings.ToUpper(m))
		fmt.Fprintf(bb, "#define GO_%v_H\n\n", strings.ToUpper(m))

		fmt.Fprint(bb, "#include <stdint.h>\n\n")

		fmt.Fprint(bb, "#ifdef __cplusplus\n")
		if m == parser.MOC {
			for _, c := range getSortedClassNamesForModule(m) {
				fmt.Fprintf(bb, "class %v;\n", c)
			}
		}
		fmt.Fprint(bb, "extern \"C\" {\n#endif\n\n")

		fmt.Fprintf(bb, "struct %v_PackedString { char* data; long long len; };\n", strings.Title(m))
		fmt.Fprintf(bb, "struct %v_PackedList { void* data; long long len; };\n", strings.Title(m))
	}

	//body
	{
		for _, c := range getSortedClassesForModule(m) {
			cTemplate(bb, c, cppEnumHeader, cppFunctionHeader, ";\n")
		}
	}

	//footer
	{
		fmt.Fprint(bb, "\n#ifdef __cplusplus\n}\n#endif\n\n#endif")
	}

	return bb.Bytes()
}
