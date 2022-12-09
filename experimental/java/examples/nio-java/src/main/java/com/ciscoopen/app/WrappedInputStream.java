package com.ciscoopen.app;

import nasp.Connection;

import java.io.IOException;
import java.io.InputStream;

public class WrappedInputStream extends InputStream {
    Connection conn;

    public WrappedInputStream(Connection conn) {
        this.conn = conn;
    }

    @Override
    public int read() throws IOException {
        byte[] buff = new byte[1];
        return read(buff, 0, 0);
    }
    @Override
    public int read(byte[] b, int off, int len) throws IOException {
        try {
            if (b.length < len ) {
                throw new IOException("Buffer size is not big enough");
            } else if (b.length == len ){
                return (int)this.conn.read(b);
            } else {
                byte[] buffer = new byte[len];
                System.arraycopy(b, 0, buffer, 0, len);
                return (int)this.conn.read(buffer);
            }
        } catch (Exception e) {
            if ("EOF".equals(e.getMessage())) {
                return -1;
            }
            throw new IOException(e);
        }
    }

}
