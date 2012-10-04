/*-
 * Copyright (c) 2012 Caoimhe Chaos <caoimhechaos@protonmail.com>,
 *                    Ancient Solutions. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions  of source code must retain  the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions  in   binary  form  must   reproduce  the  above
 *    copyright  notice, this  list  of conditions  and the  following
 *    disclaimer in the  documentation and/or other materials provided
 *    with the distribution.
 *
 * THIS  SOFTWARE IS  PROVIDED BY  ANCIENT SOLUTIONS  AND CONTRIBUTORS
 * ``AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO,  THE IMPLIED WARRANTIES OF  MERCHANTABILITY AND FITNESS
 * FOR A  PARTICULAR PURPOSE  ARE DISCLAIMED.  IN  NO EVENT  SHALL THE
 * FOUNDATION  OR CONTRIBUTORS  BE  LIABLE FOR  ANY DIRECT,  INDIRECT,
 * INCIDENTAL,   SPECIAL,    EXEMPLARY,   OR   CONSEQUENTIAL   DAMAGES
 * (INCLUDING, BUT NOT LIMITED  TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE,  DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT  LIABILITY,  OR  TORT  (INCLUDING NEGLIGENCE  OR  OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package urlconnection

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"time"

	"github.com/4ad/doozer"
)

type doozerConnection struct {
	doozer_conn *doozer.Conn
}

/**
 * Set the Doozer configuration parameters to be used.
 */
func SetupDoozer(buri, uri string) error {
	var err error
	var doozer_conn *doozer.Conn

	doozer_conn, err = doozer.DialUri(uri, buri)
	if err != nil {
		RegisterConnectionHandler("dz", doozerConnection{
			doozer_conn: doozer_conn,
		})
	}
	return err
}

/**
 * Connect to a host:port pair given in a Doozer file.
 * Makes a TCP connection to the given host:port pair.
 */
func (conn doozerConnection) Connect(dest *url.URL) (net.Conn, error) {
	var info *doozer.FileInfo
	var data []byte
	var rev int64
	var err error

	if conn.doozer_conn == nil {
		return nil, errors.New("Please use SetupDoozer first")
	}

	// Query the paths at the latest revision.
	rev, err = conn.doozer_conn.Rev()
	if err != nil {
		return nil, err
	}

	info, err = conn.doozer_conn.Statinfo(rev, dest.Path)
	if err != nil {
		return nil, err
	}

	if info.IsDir {
		var names []string
		var name string
		var selected int

		names, err = conn.doozer_conn.Getdir(dest.Path, rev, 0, -1)
		if err != nil {
			return nil, err
		}

		selected = rand.Intn(len(names))
		name = fmt.Sprintf("%s/%s", dest.Path, names[selected])
		data, _, err = conn.doozer_conn.Get(name, &rev)
		if err != nil {
			return nil, err
		}
	} else {
		data, _, err = conn.doozer_conn.Get(dest.Path, &rev)
		if err != nil {
			return nil, err
		}
	}

	return net.Dial("tcp", string(data))
}

/**
 * Connect to a host:port pair given in a Doozer file.
 * Makes a TCP connection to the given host:port pair.
 * The attempt is aborted after "timeout".
 */
func (conn doozerConnection) ConnectTimeout(dest *url.URL,
	timeout time.Duration) (net.Conn, error) {
	var info *doozer.FileInfo
	var data []byte
	var rev int64
	var err error

	if conn.doozer_conn == nil {
		return nil, errors.New("Please use SetupDoozer first")
	}

	// Query the paths at the latest revision.
	rev, err = conn.doozer_conn.Rev()
	if err != nil {
		return nil, err
	}

	info, err = conn.doozer_conn.Statinfo(rev, dest.Path)
	if err != nil {
		return nil, err
	}

	if info.IsDir {
		var names []string
		var name string
		var selected int

		names, err = conn.doozer_conn.Getdir(dest.Path, rev, 0, -1)
		if err != nil {
			return nil, err
		}

		selected = rand.Intn(len(names))
		name = fmt.Sprintf("%s/%s", dest.Path, names[selected])
		data, _, err = conn.doozer_conn.Get(name, &rev)
		if err != nil {
			return nil, err
		}
	} else {
		data, _, err = conn.doozer_conn.Get(dest.Path, &rev)
		if err != nil {
			return nil, err
		}
	}

	return net.DialTimeout("tcp", string(data), timeout)
}
