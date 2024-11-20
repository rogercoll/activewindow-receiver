package x11provider

import (
	"context"
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/rogercoll/activewindowreceiver/internal/provider"
)

type x11Provider struct {
	connection *xgb.Conn
}

var _ provider.ActiveWindowProvider = (*x11Provider)(nil)

func newX11Provider(cfg *Config) (*x11Provider, error) {
	X, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}
	return &x11Provider{
		connection: X,
	}, nil
}

func (x *x11Provider) ActiveWindow(context.Context) (string, string, error) {
	// Get the window id of the root window.
	setup := xproto.Setup(x.connection)
	root := setup.DefaultScreen(x.connection).Root

	// Get the atom id (i.e., intern an atom) of "_NET_ACTIVE_WINDOW".
	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(x.connection, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		return "", "", err
	}

	// Get the atom id (i.e., intern an atom) of "_NET_WM_NAME".
	aname = "_NET_WM_NAME"
	nameAtom, err := xproto.InternAtom(x.connection, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		return "", "", err
	}

	// Get the actual value of _NET_ACTIVE_WINDOW.
	// Note that 'reply.Value' is just a slice of bytes, so we use an
	// XGB helper function, 'Get32', to pull an unsigned 32-bit integer out
	// of the byte slice. We then convert it to an X resource id so it can
	// be used to get the name of the window in the next GetProperty request.
	reply, err := xproto.GetProperty(x.connection, false, root, activeAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return "", "", err
	}
	windowId := xproto.Window(xgb.Get32(reply.Value))

	// Now get the value of _NET_WM_NAME for the active window.
	// Note that this time, we simply convert the resulting byte slice,
	// reply.Value, to a string.
	reply, err = xproto.GetProperty(x.connection, false, windowId, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return "", "", err
	}
	// fmt.Printf("Active window name: %s\n", string(reply.Value))

	return fmt.Sprintf("%X", windowId), string(reply.Value), nil
}
