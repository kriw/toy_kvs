#include <cstring>
#include <string>
#include <cstdint>
#include <unistd.h>
#include <pthread.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in.h>

namespace TCPServer {

    void setup(std::string addr, uint16_t port) {
        int sockfd = socket(AF_INET,SOCK_STREAM,0);
        struct sockaddr_in addr_in;
        addr_in.sin_family = AF_INET;
        addr_in.sin_port = htons(port);
        inet_pton(AF_INET, addr.c_str(), &addr_in.sin_addr);

        bind(sockfd, (struct sockaddr *)&addr_in, sizeof(addr_in));
        listen(sockfd, 128);
    }

    class Conn {
        int sockfd;
        const uint16_t MAX_PACKET_SIZE = 1024;
        Conn(int _sockfd) {
            sockfd = _sockfd;
        }
        ~Conn() {
            close(sockfd);
        }
        void send(char *buf, int size) {
            while(size > 0) {
                write(sockfd, buf, MAX_PACKET_SIZE);
                size -= MAX_PACKET_SIZE;
            }
        }
        void recv(char *buf) {
            //TODO use one of these functions recv, recvfrom, recvmsg
            read(sockfd, buf, MAX_PACKET_SIZE);
        }
    };
}
