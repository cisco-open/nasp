package com.ciscoopen.app;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.net.SocketAddress;
import java.net.SocketException;

import nasp.Connection;

public class WrappedSocket extends Socket {
    Connection conn;
    public WrappedSocket(Connection conn) {
        this.conn = conn;
    }

    @Override
    public InputStream getInputStream() throws IOException {
        return new WrappedInputStream(this.conn);
    }

    @Override
    public OutputStream getOutputStream() throws IOException {
        return new WrappedOutputStream(this.conn);
    }

    @Override
    public void bind(SocketAddress bindPoint) throws IOException {

    }

    public void setReuseAddress(boolean on) throws SocketException {

    }
}
