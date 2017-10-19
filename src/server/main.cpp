#include <string>
#include <cstdint>

namespace KVS {
    const std::string FILE_DIR = "./files";

namespace detail {
    void saveFile();
    void getFile();
    void setFile();
    void registerFiles();

    void backgroudRead();
    void requestHandler();
    void handleRequest();
};

    void Serve(std::string addr, uint16_t port) {
        //run server
        
        //compile yara rules

        //run malware scanner

        while(true) {
            //listen
        }
    };
};
