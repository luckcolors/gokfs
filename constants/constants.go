package constants

const R uint64 = 160                       //Number of bits in Reference ID
const C uint64 = 131072                    //Number of bytes in a file chunk
const S uint64 = 32 * (1024 * 1024 * 1024) //Number of bytes in a Sbucket
const B uint64 = 2 ^ 8                     //Number of columns in a Btable
const D uint64 = 8                         //Number of bits for the distance calculation
const HASH string = "ripemd160"            //OpenSSL id for key hashing algorithm
const SBUCKET_IDLE uint64 = 60000          //Time to wait before idle event
