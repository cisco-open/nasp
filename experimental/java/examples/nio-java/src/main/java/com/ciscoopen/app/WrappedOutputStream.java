package com.ciscoopen.app;

import nasp.Connection;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.ByteBuffer;

public class WrappedOutputStream extends OutputStream {
    Connection conn;

    public WrappedOutputStream(Connection conn) {
        this.conn = conn;
    }
    
    @Override
    public void write(int b) throws IOException {
        final byte[] buffer = new byte[]{(byte)b};
        write(buffer, 0, 1);
    }

    @Override
    public void write(byte[] b, int off, int len) throws IOException {
        try {
            byte[] buffer = new byte[b.length];
            if (len >= 0) System.arraycopy(b, 0, buffer, 0, len);
            this.conn.write(buffer);
        } catch (Exception e) {
            throw new IOException("could not write to nasp socket");
        }
    }
    
}
