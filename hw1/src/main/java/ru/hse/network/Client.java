package ru.hse.network;

import java.io.*;
import java.net.*;
import java.util.Random;

public class Client {
    public static void main(String[] args) {
        if (args.length != 5) {
            System.out.println("Использование: java Client <IP> <port> <N> <M> <Q>");
            return;
        }

        String serverIP = args[0];
        int port = Integer.parseInt(args[1]);
        int N = Integer.parseInt(args[2]);
        int M = Integer.parseInt(args[3]);
        int Q = Integer.parseInt(args[4]);

        try (Socket socket = new Socket(serverIP, port)) {
            socket.setTcpNoDelay(true);
            
            DataOutputStream out = new DataOutputStream(socket.getOutputStream());
            DataInputStream in = new DataInputStream(socket.getInputStream());
            
            Random random = new Random();
            
            System.out.println("bytes,average_time_ms");
            
            for (int k = 0; k < M; k++) {
                int arraySize = N * k + 8;
                long totalTime = 0;
                
                for (int q = 0; q < Q; q++) {
                    System.err.printf(
                            "progress: k=%d/%d, q=%d/%d, bytes=%d%n",
                            k + 1, M,
                            q + 1, Q,
                            arraySize
                    );
                    byte[] data = new byte[arraySize];
                    random.nextBytes(data);
                    
                    long startTime = System.currentTimeMillis();
                    
                    out.writeInt(arraySize);
                    out.write(data);
                    out.flush();
                    
                    in.readUTF();
                    
                    long endTime = System.currentTimeMillis();
                    totalTime += (endTime - startTime);
                }
                
                double avgTime = (double) totalTime / Q;
                System.out.println(arraySize + "," + avgTime);
            }
            
            out.writeInt(0);
            out.flush();
            
        } catch (IOException e) {
            System.err.println("ошибка соединения " + e.getMessage());
            e.printStackTrace();
        }
    }
}
