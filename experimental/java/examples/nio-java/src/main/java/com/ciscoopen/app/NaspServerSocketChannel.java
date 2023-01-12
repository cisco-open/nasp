package com.ciscoopen.app;

import nasp.TCPListener;

import java.io.FileDescriptor;
import java.io.IOException;
import java.net.ServerSocket;
import java.nio.channels.*;
import java.nio.channels.spi.SelectorProvider;

public class NaspServerSocketChannel {
    private final SelectableChannel channel;
    private WrappedServerSocket socket;

    public NaspServerSocketChannel(SelectableChannel channel) {
        this.channel = channel;
    }

    public WrappedServerSocket socket() {
        try {
            if (this.socket == null) {
                this.socket = new WrappedServerSocket(10000);
            }
            return socket;
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    public final SelectableChannel configureBlocking(boolean block) throws IOException {
        return this.channel.configureBlocking(block);
    }

    public final SelectionKey register(Selector sel, int ops, Object att) throws ClosedChannelException {
        return this.channel.register(sel, ops, att);
    }

    public final SelectionKey register(Selector sel, int ops)
            throws ClosedChannelException
    {
        return register(sel, ops, null);
    }

    public SocketChannel accept() throws IOException {
        return this.socket.accept().getChannel();
    }

}
