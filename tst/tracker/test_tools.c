#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>
#include "tools.h"

void test_announce() {
    printf(">>> announce...");

    // 1 file
    announceData data = announceCheck(
            "announce listen 2222 seed [teemo 14 7 jzichfnt8SBYA8NS8AZNY8SN9kzox83h]");

    File *expected_files = malloc(sizeof(File));
    expected_files->name = strdup("teemo");
    expected_files->size = 14;
    expected_files->pieceSize = 7;
    strncpy(expected_files->key, "jzichfnt8SBYA8NS8AZNY8SN9kzox83h", 33);

    announceData expected_data = {.port=2222, .nb_files=1, .files = expected_files, .nb_leech_keys = 0};//, .files};
    assert(announceStructCmp(data, expected_data));

    //2 files
    announceData data2 = announceCheck(
            "announce listen 2522 seed [ferIV 120 13 9kOz8SB18SBYA8NS8AZNY8SN9kzox83h interruption 94 3 8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h]");
    //printAnnounceData(data);
    //printAnnounceData(data2);
    free_announceData(&data);
    free_announceData(&data2);
    free_regex(announce_regex());

    free(expected_files);

    printf("\033[92mpassed\033[39m\n");
}

void test_look() {
    lookData data3 = lookCheck("look [filename=\"Enfin\"]");
    lookData data4 = lookCheck("look [filesize=\"9128\" filename=\"Alttab\"]");
    //printLookData(data3);
    //printLookData(data4);
    free_regex(look_regex());
    free_lookData(&data3);
    free_lookData(&data4);

}

void test_getfile() {
    getfileData getfile = getfileCheck("getfile jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h");
    //printGetFileData(getfile);
    free_regex(comparison_regex());
}

int main() {
    test_announce();
    test_look();
    test_getfile();
    return 0;
}