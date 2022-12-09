package com.ciscoopen.app;

import java.nio.channels.ServerSocketChannel;
import java.nio.channels.spi.SelectorProvider;

public abstract class WrappedServerSocketChannel extends ServerSocketChannel {
    /**
     * Initializes a new instance of this class.
     *
     * @param provider The provider that created this channel
     */
    protected WrappedServerSocketChannel(SelectorProvider provider) {
        super(provider);
    }

    @Override
    public WrappedServerSocket socket() {
        return null;
    }

}
