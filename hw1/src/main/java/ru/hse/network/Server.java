package ru.hse.network;

import java.io.*;
import java.net.*;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

public class Server {
    private static final int MAX_PAYLOAD_BYTES = 8 * 1024 * 1024; // 8 MB

    public static void main(String[] args) {
        if (args.length != 1) {
            System.out.println("Использование: java Server <port>");
            return;
        }

        int port = Integer.parseInt(args[0]);
        
        try (ServerSocket serverSocket = new ServerSocket(port)) {
            System.out.println("сервер поднялся на порту " + port);
            
            while (true) {
                try (Socket clientSocket = serverSocket.accept()) {
                    System.out.println("подключился клиент: " + clientSocket.getInetAddress());
                    handleClient(clientSocket);
                } catch (IOException e) {
                    System.err.println("ошибка клиента " + e.getMessage());
                }
            }
        } catch (IOException e) {
            System.err.println("занят порт " + port);
            e.printStackTrace();
        }
    }

    private static void handleClient(Socket socket) throws IOException {
        InputStream in = socket.getInputStream();
        OutputStream out = socket.getOutputStream();
        DataInputStream dataIn = new DataInputStream(in);
        DataOutputStream dataOut = new DataOutputStream(out);
        
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy.MM.dd HH:mm:ss");
        
        try {
            while (true) {
                int length = dataIn.readInt();
                
                if (length <= 0) {
                    break;
                }

                if (length > MAX_PAYLOAD_BYTES) {
                    System.out.println("кто то прислал слишком большой пакет " + length + " байт. соединение закрыто");
                    break;
                }
                
                byte[] buffer = new byte[length];
                dataIn.readFully(buffer);
                
                String timestamp = LocalDateTime.now().format(formatter);
                dataOut.writeUTF(timestamp);
                dataOut.flush();
            }
        } catch (EOFException e) {
            System.out.println("клиент отключен");
        }
    }
}
