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

    announceData expected_data = {.port=2222, .nb_files=1, .files = expected_files, .nb_leech_keys = 0, .is_valid=1};
    assert(announceStructCmp(data, expected_data));

    // 1 leech key
    announceData data2 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnA8NS8AZNY8SN9kzox83h]");
    char *expected_leech_key = "jzichfnt8SBnA8NS8AZNY8SN9kzox83h";
    announceData expected_data2 = {.port=4222, .nb_files=0, .nb_leech_keys = 1, .leechKeys = &expected_leech_key, .is_valid=1};
    assert(announceStructCmp(data2, expected_data2));

    //2 leech keys
    announceData data3 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnA8N78AZNY8SN9kzox83h jzichfnt818nA8N78AZNY8SN9kzox83h]");
    char *expected_leech_key1 = "jzichfnt8SBnA8N78AZNY8SN9kzox83h";
    char *expected_leech_key2 = "jzichfnt818nA8N78AZNY8SN9kzox83h";
    char *expected_leech_keys[2];
    expected_leech_keys[0] = expected_leech_key1;
    expected_leech_keys[1] = expected_leech_key2;
    announceData expected_data3 = {.port=4222, .nb_files=0, .nb_leech_keys = 2, .leechKeys = expected_leech_keys, .is_valid=1};
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
    announceData expected_data4 = {.port=2522, .nb_files=2, .files=expected_files4,.nb_leech_keys = 2, .leechKeys = expected_leech_keys4, .is_valid=1};
    assert(announceStructCmp(data4, expected_data4));

    // Not valid
    announceData data5 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnANS8AZNY8SN9kzox83h]");
    announceData not_valid = {.is_valid=0};
    assert(!announceStructCmp(data5, not_valid));
    announceData data6 = announceCheck(
            "announce listen 4222 seed [fgh 18 9 s]");
    assert(!announceStructCmp(data6, not_valid));
    announceData data7 = announceCheck(
            "announce listen 4222 seed [fgh 18 9i jzichint8SBnANS8AZNsY8SN9kzox83h]");
    assert(!announceStructCmp(data7, not_valid));

    announceData data8 = announceCheck(
            "announce listen 4222 seed [fgh 18o 2 jzichint8SBnANS8AZNsY8SN9kzox83h]");
    assert(!announceStructCmp(data8, not_valid));
    announceData data9 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnANS8AZNYsl8SN9kzox83h]");
    assert(!announceStructCmp(data9, not_valid));


    free_announceData(&data);
    free_announceData(&data2);
    free_announceData(&data3);
    free_announceData(&data4);
    free_announceData(&data5);
    free_announceData(&data6);
    free_announceData(&data7);
    free_announceData(&data8);
    free_announceData(&data9);
    free_regex(announce_regex());
    free(expected_files);

    printf("\033[92mpassed\033[39m\n");
}

void test_look() {
    printf(">>> look...");
    lookData data = lookCheck("look [filename=\"Enfin\"]");
    lookData data2 = lookCheck("look [filesize>=\"9128\" filename=\"Alttab\"]");
    char *name = "Enfin";
    criterion crit = {.criteria=FILENAME, .op=EQ, .value_type=STR,.value.value_str = name};
    lookData expected_look_data = {.nb_criterions=1, .criterions=&crit, .is_valid=1};
    assert(lookStructCmp(data, expected_look_data));

    char *name2 = "Alttab";
    criterion crit2 = {.criteria=FILENAME, .op=EQ, .value_type=STR,.value.value_str = name2};
    criterion crit3 = {.criteria=FILESIZE, .op=GE, .value_type=INT, .value.value_int = 9128};
    criterion crits[2];
    crits[0] = crit2;
    crits[1] = crit3;
    lookData expected_look_data2 = {.nb_criterions=2, .criterions=crits, .is_valid=1};
    assert(lookStructCmp(data2, expected_look_data2));

    free_regex(look_regex());
    free_lookData(&data);
    free_lookData(&data2);

    printf("\033[92mpassed\033[39m\n");
}

void test_getfile() {
    printf(">>> getfile...");
    getfileData data = getfileCheck("getfile jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h");
    getfileData expected_data = {.key="jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h", .is_valid=1};
    assert(getfileStructCmp(data, expected_data));

    // Not valid
    getfileData not_valid = {.is_valid=0};
    getfileData data2 = getfileCheck("getfile jzicsfnt8SBA8NS8AZNY8SN9dkzo83h");
    assert(!getfileStructCmp(data2, not_valid));
    getfileData data3 = getfileCheck("getfile jzi784sfnt8SBA8NS8AZNY8SN9dkzo83h");
    assert(!getfileStructCmp(data3, not_valid));
    free_regex(comparison_regex());
    printf("\033[92mpassed\033[39m\n");
}

int main() {
    test_announce();
    test_look();
    test_getfile();
    return 0;
}