// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"io"
	"os"
)

// logWriter specifies the output target of the logTraceWriter; replaced in tests.
var logWriter io.Writer = os.Stdout

// setLogWriter sets the io.Writer that any new logTraceWriter will write to and returns a function
// which will return the io.Writer to its original value.
func setLogWriter(w io.Writer) func() {
	tmp := logWriter
	logWriter = w
	return func() { logWriter = tmp }
}
