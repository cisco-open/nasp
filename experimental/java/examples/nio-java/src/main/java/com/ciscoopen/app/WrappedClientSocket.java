package com.ciscoopen.app;

import nasp.Connection;
import nasp.Nasp;
import nasp.TCPDialer;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.net.SocketAddress;

public abstract class WrappedClientSocket extends Socket {
    TCPDialer dialer;
    Connection conn;
    public WrappedClientSocket() throws IOException {
        try {
            dialer = Nasp.newTCPDialer("https://localhost:16443/config",
                    "test-tcp-16362813-F46B-41AC-B191-A390DB1F6BDF",
                    "16362813-F46B-41AC-B191-A390DB1F6BDF");
        } catch (Exception e) {
            throw new IOException("could not get nasp tcp dialer");
        }
    }

    public void connect(SocketAddress endpoint) throws IOException {
        try {
            conn = dialer.dial();
        } catch (Exception e) {
            throw new IOException("could not get nasp tcp dialer");
        }
    }

    @Override
    public InputStream getInputStream() {
        return new WrappedInputStream(this.conn);
    }

    @Override
    public OutputStream getOutputStream() {
        return new WrappedOutputStream(this.conn);
    }
}
