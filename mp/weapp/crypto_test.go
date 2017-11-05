// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package weapp

import (
	"testing"

	"github.com/issue9/assert"
)

func TestDecode(t *testing.T) {
	a := assert.New(t)

	appid := "wx4f4bc4dec97d474b"
	sessionKey := "tiihtNczf5v6AKRyjwEUhQ=="
	iv := "r7BXXKkLb8qrSNn05n0qiA=="
	encryptedData := `CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZM
QmRzooG2xrDcvSnxIMXFufNstNGTyaGS
9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+
3hVbJSRgv+4lGOETKUQz6OYStslQ142d
NCuabNPGBzlooOmB231qMM85d2/fV6Ch
evvXvQP8Hkue1poOFtnEtpyxVLW1zAo6
/1Xx1COxFvrc2d7UL/lmHInNlxuacJXw
u0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn
/Hz7saL8xz+W//FRAUid1OksQaQx4CMs
8LOddcQhULW4ucetDf96JcR3g0gfRK4P
C7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB
6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns
/8wR2SiRS7MNACwTyrGvt9ts8p12PKFd
lqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYV
oKlaRv85IfVunYzO0IKXsyl7JCUjCpoG
20f0a04COwfneQAGGwd5oa+T8yO5hzuy
Db/XcxxmK01EpqOyuxINew==`

	data, watermark, err := Decode(appid, sessionKey, encryptedData, iv)
	a.NotError(err).NotNil(data).NotNil(watermark)
	a.Equal(appid, watermark.Appid)
	a.TB().Log(string(data))
}
