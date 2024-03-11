// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.

package namingschema

import "github.com/DataDog/dd-trace-go/v2/v1internal/namingschema"

func ServiceName(fallback string) string {
	return namingschema.ServiceName(fallback)
}

func ServiceNameOverrideV0(fallback, overrideV0 string) string {
	return namingschema.ServiceNameOverrideV0(fallback, overrideV0)
}
