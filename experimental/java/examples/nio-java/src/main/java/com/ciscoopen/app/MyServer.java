package com.ciscoopen.app;

import java.io.*;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.*;
import java.util.Iterator;
import java.util.Set;

public class MyServer {
    public static void main(String[] args) {
        initializeServer();
//        initializeClient();
    }
    public static void initializeClient() {
        try {
            SocketChannel channel = SocketChannel.open();
            channel.connect(new InetSocketAddress("localhost", 10000));
            System.out.println("Connection Set: " + channel.getRemoteAddress());
            ByteBuffer buffer = ByteBuffer.allocate(256);
            buffer.clear();
            String newData = "This is string from the client";
            buffer.put(newData.getBytes());
            buffer.flip();
            while (buffer.hasRemaining()) {
                channel.write(buffer);
            }
            while (channel.read(buffer) < 0) {
                System.out.println(buffer);
                buffer.clear();
            }
            channel.close();

        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static void initializeServer() {
        try {
            NaspServerSocketChannel server = new NaspServerSocketChannel(ServerSocketChannel.open());
//            ServerSocketChannel server = ServerSocketChannel.open();
            server.socket().bind(new InetSocketAddress("localhost", 10000));
            server.socket().setReuseAddress(true);
            server.configureBlocking(false);

            Selector selector = Selector.open();
            server.register(selector, SelectionKey.OP_ACCEPT);

            ByteBuffer buffer = ByteBuffer.allocate(256);
            while (true) {
                int channelCount = selector.select(10000);
                System.out.println("After select");
                if (channelCount > 0) {
                    Set<SelectionKey> keys = selector.selectedKeys();
                    Iterator<SelectionKey> iterator = keys.iterator();
                    while (iterator.hasNext()) {
                        SelectionKey key = iterator.next();
                        iterator.remove();
                        if (key.isAcceptable()) {
                            SocketChannel client = server.accept();
                            client.configureBlocking(false);
                            client.register(selector, SelectionKey.OP_READ, client.socket().getPort());
                        } else if (key.isReadable()) {
                            SocketChannel client = (SocketChannel) key.channel();
                            p("port: " + key.attachment());
                            if (client.read(buffer) < 0) {
                                key.cancel();
                                client.close();
                            } else {
                                buffer.flip(); // read from the buffer
                                client.write(buffer);
                                buffer.clear();
                            }
                        }
                    }
                }
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    private static void p(String s) {
        System.out.println(s);
    }
}