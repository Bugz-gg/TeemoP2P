#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>
#include "tools.h"

void test_announce() {
    printf(">>> announce...");

    // 1 file
    announceData data = announceCheck(
            "announce listen 2222 seed [teemo 14 7 jzichfnt8SBYA8NS8AZNY8SN9kzox83h]\n");
    printAnnounceData(data);

    File *expected_files = malloc(sizeof(File));
    strcpy(expected_files->name, "teemo");
    expected_files->size = 14;
    expected_files->pieceSize = 7;
    strncpy(expected_files->key, "jzichfnt8SBYA8NS8AZNY8SN9kzox83h", 33);
    announceData expected_data = {.port=2222, .nb_files=1, .files = expected_files, .nb_leech_keys = 0, .is_valid=1};
    assert(announceStructCmp(data, expected_data));

    // 1 leech key
    announceData data2 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnA8NS8AZNY8SN9kzox83h]\r\n");
    char *expected_leech_key = "jzichfnt8SBnA8NS8AZNY8SN9kzox83h";
    announceData expected_data2 = {.port=4222, .nb_files=0, .nb_leech_keys = 1, .leechKeys = &expected_leech_key, .is_valid=1};
    assert(announceStructCmp(data2, expected_data2));

    //2 leech keys
    announceData data3 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnA8N78AZNY8SN9kzox83h jzichfnt818nA8N78AZNY8SN9kzox83h]\n");
    char *expected_leech_key1 = "jzichfnt8SBnA8N78AZNY8SN9kzox83h";
    char *expected_leech_key2 = "jzichfnt818nA8N78AZNY8SN9kzox83h";
    char *expected_leech_keys[2];
    expected_leech_keys[0] = expected_leech_key1;
    expected_leech_keys[1] = expected_leech_key2;
    announceData expected_data3 = {.port=4222, .nb_files=0, .nb_leech_keys = 2, .leechKeys = expected_leech_keys, .is_valid=1};
    assert(announceStructCmp(data3, expected_data3));

    //2 files and 2 leech keys
    announceData data4 = announceCheck(
            "announce listen 2522 seed [ferIV 120 12 9kOz8SB18SBYA8NS8AZNY8SN9kzox83h interruption 94 2 8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h] leech [jzichfnt818nA8N78jzkY8SN9kzox83h jzichfnt818nA8N7rlpNY8SN9kzox83h]\r\n");
    File *expected_files4 = malloc(2 * sizeof(File));
    strcpy(expected_files4[0].name, "ferIV");
    expected_files4[0].size = 120;
    expected_files4[0].pieceSize = 12;
    strncpy(expected_files4[0].key, "9kOz8SB18SBYA8NS8AZNY8SN9kzox83h", 33);
    strcpy(expected_files4[1].name, "interruption");
    expected_files4[1].size = 94;
    expected_files4[1].pieceSize = 2;
    strncpy(expected_files4[1].key, "8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h", 33);
    char *expected_leech_key1_4 = "jzichfnt818nA8N78jzkY8SN9kzox83h";
    char *expected_leech_key2_4 = "jzichfnt818nA8N7rlpNY8SN9kzox83h";
    char *expected_leech_keys4[2];
    expected_leech_keys4[0] = expected_leech_key1_4;
    expected_leech_keys4[1] = expected_leech_key2_4;
    announceData expected_data4 = {.port=2522, .nb_files=2, .files=expected_files4, .nb_leech_keys = 2, .leechKeys = expected_leech_keys4, .is_valid=1};
    assert(announceStructCmp(data4, expected_data4));

    // Not valid
    announceData data5 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnANS8AZNY8SN9kzox83h]\r\n");
    announceData not_valid = {.is_valid=0, .nb_files=0, .nb_leech_keys=0, .port=0};
    assert(announceStructCmp(data5, not_valid));
    announceData data6 = announceCheck(
            "announce listen 4222 seed [fgh 18 jzichfnt8lBnANS8AZNY8SN9kzox83h]\n");
    assert(announceStructCmp(data6, not_valid));
    announceData data7 = announceCheck(
            "announce listen 4222 seed [fgh 18 9i jzichint8SBnANS8AZNsY8SN9kzox83h]\r\n");
    assert(announceStructCmp(data7, not_valid));

    announceData data8 = announceCheck(
            "announce listen 4222 seed [fgh 18o 2 jzichint8SBnANS8AZNsY8SN9kzox83h]\n");
    assert(announceStructCmp(data8, not_valid));
    announceData data9 = announceCheck(
            "announce listen 4222 seed [] leech [jzichfnt8SBnANS8AZNYsl8SN9kzox83h]\r\n");
    assert(announceStructCmp(data9, not_valid));


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
    free(expected_data4.files);
    free(expected_files);

    printf("\033[92mpassed\033[39m\n");
}

void test_look() {
    printf(">>> look...");
    lookData data = lookCheck("look [filename=\"Enfin\"]\n");
    char *name = "Enfin";
    criterion crit = {.criteria=FILENAME, .op=EQ, .value_type=STR, .value.value_str = name};
    lookData expected_look_data = {.nb_criterions=1, .criterions=&crit, .is_valid=1};
    assert(lookStructCmp(data, expected_look_data));

    lookData data2 = lookCheck("look [filesize>=\"9128\" filename=\"Alttab\"]\n");
    char *name2 = "Alttab";
    criterion crit2 = {.criteria=FILENAME, .op=EQ, .value_type=STR, .value.value_str = name2};
    criterion crit3 = {.criteria=FILESIZE, .op=GE, .value_type=INT, .value.value_int = 9128};
    criterion crits[2];
    crits[0] = crit2;
    crits[1] = crit3;
    lookData expected_look_data2 = {.nb_criterions=2, .criterions=crits, .is_valid=1};
    assert(lookStructCmp(data2, expected_look_data2));

    // Not valid
    lookData data3 = lookCheck("look [filesie=\"9128\" filename=m\"Alltab\"]\n");
    lookData not_valid = {.nb_criterions=0, .is_valid=0};
    assert(lookStructCmp(data3, not_valid));
    lookData data4 = lookCheck("look [filesie\"9128\" filename=\"Alltab\"]\r\n");
    assert(lookStructCmp(data4, not_valid));
    lookData data5 = lookCheck("look [filesie>=\"9128\" filename=\"Alltab\"]\r\n");
    assert(lookStructCmp(data5, not_valid));

    lookData data6 = lookCheck("look [filesize>=\"9128\" filename=\"Alttab\" key=\"dsqklmdsqklmdsqklmdsqklmdsqklmnj\"]\r\n");
    criterion crit4 = {.criteria=KEY, .op=EQ, .value_type=STR, .value.value_str = "dsqklmdsqklmdsqklmdsqklmdsqklmnj"};
    criterion crits2[3];
    crits2[0] = crit2;
    crits2[1] = crit3;
    crits2[2] = crit4;
    lookData expected_look_data6 = {.nb_criterions=3, .criterions=crits2, .is_valid=1};
    assert(lookStructCmp(data6, expected_look_data6));

    free_regex(look_regex());
    free_lookData(&data);
    free_lookData(&data2);
    free_lookData(&data3);
    free_lookData(&data4);

    printf("\033[92mpassed\033[39m\n");
}

void test_getfile() {
    printf(">>> getfile...");
    getfileData data = getfileCheck("getfile jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h\r\n");
    getfileData expected_data = {.key="jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h", .is_valid=1};
    assert(getfileStructCmp(data, expected_data));

    // Not valid
    getfileData not_valid = {.is_valid=0};
    getfileData data2 = getfileCheck("getfile jzicsfnt8SBA8NS8AZNY8SN9dkzo83h\r\n");
    assert(getfileStructCmp(data2, not_valid));
    getfileData data3 = getfileCheck("getfile jzi784sfnt8SBA8NS8AZNY8SN9dkzo83h\r\n");
    assert(getfileStructCmp(data3, not_valid));
    free_regex(comparison_regex());
    printf("\033[92mpassed\033[39m\n");
}

void test_update() {
    printf(">>> update...");
    updateData data = updateCheck("update seed [jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h]\r\n");
    char *key = "jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h";
    updateData expected_data = {.nb_keys=1, .keys=&key, .nb_leech=0, .is_valid=1};
    assert(updateStructCmp(data, expected_data));

    char *leech = "jzicsfnt8SBYA8NerAZNY8SN9dkzo83h";
    updateData data2 = updateCheck("update seed [] leech [jzicsfnt8SBYA8NerAZNY8SN9dkzo83h]\n");
    updateData expected_data2 = {.nb_keys=0, .leech=&leech, .nb_leech=1, .is_valid=1};
    assert(updateStructCmp(data2, expected_data2));

    char *key1 = "jzicsfnt8SBYA8NerAZNl8SN9dkzo83h";
    char *key2 = "jzicsfnlzbBYA8NerAZNY8SN9dkzo83h";
    char *keys[2] = {key1, key2};
    updateData data3 = updateCheck("update seed [jzicsfnt8SBYA8NerAZNl8SN9dkzo83h jzicsfnlzbBYA8NerAZNY8SN9dkzo83h] leech [jzicsfnt8SBYA8NerAZNY8SN9dkzo83h]\r\n");
    updateData expected_data3 = {.nb_keys=2, .keys=keys, .nb_leech=1, .leech=&leech, .is_valid=1};
    assert(updateStructCmp(data3, expected_data3));

    key1 = "jzicsfnt8SBYA8NerAZNllSN9dkzo83h";
    key2 = "jzicsfnlzbBYA8NerAZNY8pN9dkzo83h";
    keys[0] = key1;
    keys[1] = key2;
    char *leech1 = "lzicsfnt8SBYA8NerAZNllSN9dkzo83h";
    char *leech2 = "jzicsfnt8SleA8NerAZNllSN9dkzo83h";
    char *leech3 = "jzicsfnt8SBYA2NerAZNllSN9dkzo83h";
    char *leechs[3] = {leech1, leech2, leech3};
    updateData data4 = updateCheck("update seed [jzicsfnt8SBYA8NerAZNllSN9dkzo83h jzicsfnlzbBYA8NerAZNY8pN9dkzo83h] leech [lzicsfnt8SBYA8NerAZNllSN9dkzo83h jzicsfnt8SleA8NerAZNllSN9dkzo83h jzicsfnt8SBYA2NerAZNllSN9dkzo83h]\r\n");
    updateData expected_data4 = {.nb_keys=2, .keys=keys, .nb_leech=3, .leech=leechs, .is_valid=1};
    assert(updateStructCmp(data4, expected_data4));

    // Not valid
    updateData not_valid = {.is_valid=0};
    updateData data5 = updateCheck("update seed [jzicsfnt8SBA8NS8AZNY8SN9dkzo83h]\r\n");
    assert(updateStructCmp(data5, not_valid));
    updateData data6 = updateCheck("update seed [jzi784sfnt8SBA8NS8AZNY8SN9dkzo83h]\r\n");
    assert(updateStructCmp(data6, not_valid));
    free_regex(update_regex());
    free_updateData(&data5);
    free_updateData(&data6);
    free_updateData(&data4);
    free_updateData(&data3);
    free_updateData(&data2);
    free_updateData(&data);
    printf("\033[92mpassed\033[39m\n");
}

int main() {
    test_announce();
    test_look();
    test_getfile();
    test_update();
    return 0;
}