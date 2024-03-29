#include "tools.h"

int main() {
    announceData data = announceCheck(
            "announce listen 2222 seed [teemo 14 5 jzichfnt8SBYA8NS8AZNY8SN9kzox83h]");
    announceData data2 = announceCheck(
            "announce listen 2522 seed [ferIV 120 13 9kOz8SB18SBYA8NS8AZNY8SN9kzox83h interruption 94 3 8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h]");

    lookData data3 = lookCheck("look [filename=\"Enfin\"]");
    lookData data4 = lookCheck("look [filesize=\"9128\" filename=\"Alttab\"]");

    getfileData getfile = getfileCheck("getfile jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h");


    printAnnounceData(data);
    printAnnounceData(data2);
    printLookData(data3);
    printLookData(data4);
    printGetFileData(getfile);

    free_regex(announce_regex());
    free_regex(look_regex());
    free_regex(comparison_regex());
    free_announceData(&data);
    free_announceData(&data2);
    free_lookData(&data3);
    free_lookData(&data4);

    return 0;
}