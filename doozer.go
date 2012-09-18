/*-
 * Copyright (c) 2012 Tonnerre Lombard <tonnerre@ancient-solutions.com>,
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
	"github.com/4ad/doozer"
	"net"
	"net/url"
)

var doozer_conn *doozer.Conn

/**
 * Set the Doozer configuration parameters to be used.
 */
func SetupDoozer(buri, uri string) error {
	var err error

	doozer_conn, err = doozer.DialUri(uri, buri)
	return err
}

/**
 * Connect to a host:port pair given in a Doozer file.
 * Makes a TCP connection to the given host:port pair.
 */
func doozerConnect(dest *url.URL) (net.Conn, error) {
	var data []byte
	var err error

	if doozer_conn == nil {
		return nil, errors.New("Please use SetupDoozer first")
	}

	// TODO(tonnerre): Deal with directories here.
	data, _, err = doozer_conn.Get(dest.Path, nil)
	if err != nil {
		return nil, err
	}

	return net.Dial("tcp", string(data))
}

/**
 * Register the connection handler for TCP.
 */
func init() {
	RegisterConnectionHandler("dz", doozerConnect)
}