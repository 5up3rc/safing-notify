package icons

import "github.com/mattn/go-gtk/gdkpixbuf"

var (
	ByLevel = []*gdkpixbuf.Pixbuf{
		gdkpixbuf.NewPixbufFromData(Off128),
		gdkpixbuf.NewPixbufFromData(Dynamic128),
		gdkpixbuf.NewPixbufFromData(Secure128),
		gdkpixbuf.NewPixbufFromData(Fortress128),
	}

	DataByLevel = []*gdkpixbuf.PixbufData{
		&Off128,
		&Dynamic128,
		&Secure128,
		&Fortress128,
	}
)
