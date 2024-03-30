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

    // 1 leech key
    announceData data2 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnA8NS8AZNY8SN9kzox83h]");
    char *expected_leech_key = "jzichfnt8SBnA8NS8AZNY8SN9kzox83h";
    announceData expected_data2 = {.port=4222, .nb_files=0, .nb_leech_keys = 1, .leechKeys = &expected_leech_key};//, .files};
    assert(announceStructCmp(data2, expected_data2));

    //2 leech keys
    announceData data3 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnA8N78AZNY8SN9kzox83h jzichfnt818nA8N78AZNY8SN9kzox83h]");
    char *expected_leech_key1 = "jzichfnt8SBnA8N78AZNY8SN9kzox83h";
    char *expected_leech_key2 = "jzichfnt818nA8N78AZNY8SN9kzox83h";
    char *expected_leech_keys[2];
    expected_leech_keys[0] = expected_leech_key1;
    expected_leech_keys[1] = expected_leech_key2;
    announceData expected_data3 = {.port=4222, .nb_files=0, .nb_leech_keys = 2, .leechKeys = expected_leech_keys};//, .files};
    assert(announceStructCmp(data3, expected_data3));

    //2 files and 2 leech keys
    announceData data4 = announceCheck(
            "announce listen 2522 seed [ferIV 120 12 9kOz8SB18SBYA8NS8AZNY8SN9kzox83h interruption 94 2 8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h] leech [jzichfnt818nA8N78jzkY8SN9kzox83h jzichfnt818nA8N7rlpNY8SN9kzox83h]");
    File *expected_files4 = malloc(2*sizeof(File));
    expected_files4[0].name = strdup("ferIV");
    expected_files4[0].size = 120;
    expected_files4[0].pieceSize = 12;
    strncpy(expected_files4[0].key, "9kOz8SB18SBYA8NS8AZNY8SN9kzox83h", 33);
    expected_files4[1].name = strdup("interruption");
    expected_files4[1].size = 94;
    expected_files4[1].pieceSize = 2;
    strncpy(expected_files4[1].key, "8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h", 33);
    char *expected_leech_key1_4 = "jzichfnt818nA8N78jzkY8SN9kzox83h";
    char *expected_leech_key2_4 = "jzichfnt818nA8N7rlpNY8SN9kzox83h";
    char *expected_leech_keys4[2];
    expected_leech_keys4[0] = expected_leech_key1_4;
    expected_leech_keys4[1] = expected_leech_key2_4;
    announceData expected_data4 = {.port=2522, .nb_files=2, .files=expected_files4,.nb_leech_keys = 2, .leechKeys = expected_leech_keys4};//, .files};
    assert(announceStructCmp(data4, expected_data4));

    free_announceData(&data);
    free_announceData(&data2);
    free_announceData(&data3);
    free_announceData(&data4);
    free_regex(announce_regex());
    free(expected_files);

    printf("\033[92mpassed\033[39m\n");
}

void test_look() {
    printf(">>> look...");
    lookData data3 = lookCheck("look [filename=\"Enfin\"]");
    lookData data4 = lookCheck("look [filesize=\"9128\" filename=\"Alttab\"]");
    //printLookData(data3);
    //printLookData(data4);
    free_regex(look_regex());
    free_lookData(&data3);
    free_lookData(&data4);

    printf("\033[92mpassed\033[39m\n");
}

void test_getfile() {
    printf(">>> getfile...");
    getfileData getfile = getfileCheck("getfile jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h");
    //printGetFileData(getfile);
    free_regex(comparison_regex());
    (void)getfile;
    printf("\033[92mpassed\033[39m\n");
}

int main() {
    test_announce();
    test_look();
    test_getfile();
    return 0;
}