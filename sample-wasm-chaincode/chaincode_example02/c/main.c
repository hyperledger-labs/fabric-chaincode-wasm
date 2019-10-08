#define WASM_EXPORT __attribute__((visibility("default")))

/* External functions provided by wasmcc. */
extern int __print(const char *msg, int len);
extern int __get_parameter(int paramNumber, const char *result);
extern int __get_state(const char *msg, int len, const char *value);
extern int __put_state(const char *key, int keySize, const char *value, int valueSize);
extern int __delete_state(const char *msg, int len);
extern int __return_result(const char *msg, int len);

void __print_wrapper(const char *str);
void __print_wrapper_size(const char *str, int size);
void __return_result_wrapper(const char *str);
int char2int(const char *array, int n);
char *int2char(int iNumber);

/* Custom malloc as predefined malloc doesn't work */
/*
 * In case of importing starting index for memory allocation
 *
 * extern unsigned char __heap_base;
 * unsigned int bump_pointer = __heap_base;
 */

unsigned int bump_pointer = 5000;

void *malloc(unsigned long n) {
    unsigned int r = bump_pointer;
    bump_pointer += n;
    return (void *) r;
}


WASM_EXPORT
int init(int args) {

    //Should have 4 arguments
    if (args != 4) {
        const char *str = "ERROR! Incorrect number of arguments. Expecting 4";
        __print_wrapper(str);
        return -1;
    }

    //Get first parameter as 1st account's name
    const char *accountOneName = malloc(24 * sizeof(char));
    int accountOneNameLen = __get_parameter(0, accountOneName);

    //Get second parameter as 1st account's balance
    const char *accountOneBal = malloc(24 * sizeof(char));
    int accountOneBalLen = __get_parameter(1, accountOneBal);

    //Get third parameter as 2nd account's name
    const char *accountTwoName = malloc(24 * sizeof(char));
    int accountTwoNameLen = __get_parameter(2, accountTwoName);

    //Get fourth parameter as 2nd account's balance
    const char *accountTwoBal = malloc(24 * sizeof(char));
    int accountTwoBalLen = __get_parameter(3, accountTwoBal);

    //Print all values
    __print_wrapper_size(accountOneName, accountOneNameLen);
    __print_wrapper_size(accountOneBal, accountOneBalLen);
    __print_wrapper_size(accountTwoName, accountTwoNameLen);
    __print_wrapper_size(accountTwoBal, accountTwoBalLen);

    //Store account1 details
    if (__put_state(accountOneName, accountOneNameLen, accountOneBal, accountOneBalLen) == -1) {
        return -1;
    }

    //Store account2 details
    if (__put_state(accountTwoName, accountTwoNameLen, accountTwoBal, accountTwoBalLen) == -1) {
        return -1;
    }

    return 0;
}


WASM_EXPORT
int invoke(int args) {

    //Should have 3 arguments
    if (args != 3) {
        const char *str = "ERROR! Incorrect number of arguments. Expecting 3";
        __print_wrapper(str);
        return -1;
    }

    //Get first parameter as 1st account's name
    const char *accountOneName = malloc(24 * sizeof(char));
    int accountOneNameLen = __get_parameter(0, accountOneName);

    //Get second parameter as 1st account's balance
    const char *accountTwoName = malloc(24 * sizeof(char));
    int accountTwoNameLen = __get_parameter(1, accountTwoName);

    //Get third parameter as balance to transfer
    const char *amountToTransfer = malloc(24 * sizeof(char));
    int amountToTransferLen = __get_parameter(2, amountToTransfer);

    //Print all values
    __print_wrapper_size(accountOneName, accountOneNameLen);
    __print_wrapper_size(accountTwoName, accountTwoNameLen);
    __print_wrapper_size(amountToTransfer, amountToTransferLen);

    //Get balance of 1st account
    const char *accountOneBal = malloc(24 * sizeof(char));
    int accountOneBalLen = __get_state(accountOneName, accountOneNameLen, accountOneBal);


    //Convert to int
    int amountToTransferInt = char2int(amountToTransfer, amountToTransferLen);
    int accountOneBalInt = char2int(accountOneBal, accountOneBalLen);

    //Check if it have enough balance
    if (accountOneBalInt < amountToTransferInt) {
        const char *str = "amount is more than account balance";
        __return_result_wrapper(str);
        return -1;
    }

    //Get balance of 2nd account
    const char *accountTwoBal = malloc(24 * sizeof(char));
    int accountTwoBalLen = __get_state(accountTwoName, accountTwoNameLen, accountTwoBal);
    int accountTwoBalInt = char2int(accountTwoBal, accountTwoBalLen);

    //Perform credit and debit
    accountOneBalInt -= amountToTransferInt;
    accountTwoBalInt += amountToTransferInt;


    const char *str = "calling sprintf";
    __print_wrapper(str);
    /*sprintf(accountOneUpdatedBal, "%d", accountOneBalInt);
    sprintf(accountTwoUpdatedBal, "%d", accountTwoBalInt);*/

    const char *accountOneUpdatedBal = int2char(accountOneBalInt);
    const char *accountTwoUpdatedBal = int2char(accountTwoBalInt);

    const char *str2 = "sprintf success";
    __print_wrapper(str2);

    //Store account1 updated balance
    int valueSize;
    for (valueSize = 0; accountOneUpdatedBal[valueSize] != '\0'; ++valueSize);
    if (__put_state(accountOneName, accountOneNameLen, accountOneUpdatedBal, valueSize) == -1) {
        return -1;
    }

    //Store account2 updated balance
    for (valueSize = 0; accountTwoUpdatedBal[valueSize] != '\0'; ++valueSize);
    if (__put_state(accountTwoName, accountTwoNameLen, accountTwoUpdatedBal, valueSize) == -1) {
        return -1;
    }

    return 0;
}


WASM_EXPORT
int query(int args) {

    //Should have 1 argument
    if (args != 1) {
        const char *str = "ERROR! Incorrect number of arguments. Expecting name of the account to query";
        __print_wrapper(str);
        return -1;
    }

    //Get first parameter as account name
    const char *accountName = malloc(24 * sizeof(char));
    int accountNameLen = __get_parameter(0, accountName);

    //Get balance of account from state
    const char *accountBal = malloc(24 * sizeof(char));
    int accountBalLen = __get_state(accountName, accountNameLen, accountBal);

    //Return error if account balance length is negative
    if (accountBalLen == -1) {
        const char *str = "ERROR! Invalid account";
        __print_wrapper(str);
        return -1;
    }

    //Print account balance
    __print_wrapper_size(accountBal, accountBalLen);

    //Return account balance
    __return_result(accountBal, accountBalLen);

    //Return zero to show success
    return 0;
}


WASM_EXPORT
int delete(int args) {

    //Should have 1 argument
    if (args != 1) {
        const char *str = "ERROR! Incorrect number of arguments. Expecting 1";
        __print_wrapper(str);
        return -1;
    }

    //Get first parameter as account name
    const char *accountName = malloc(24 * sizeof(char));
    int accountNameLen = __get_parameter(0, accountName);

    //delete account from state
    int result = __delete_state(accountName, accountNameLen);

    //Return error if account balance length is negative
    if (result < 0) {
        const char *str = "Failed to delete state";
        __return_result_wrapper(str);
        return -1;
    }
    //Return zero to show success
    return 0;
}


void __print_wrapper_size(const char *str, int size) {
    __print(str, size);
}

void __print_wrapper(const char *str) {
    int size;
    for (size = 0; str[size] != '\0'; ++size);
    __print_wrapper_size(str, size);
}

void __return_result_wrapper(const char *str) {
    int size;
    for (size = 0; str[size] != '\0'; ++size);
    __return_result(str, size);
}

int char2int(const char *array, int n) {
    int counter = 0;
    int results = 0;
    while (1) {
        if (array[counter] == '\0') {
            break;
        } else {
            results *= 10;
            results += (int) array[counter] - 48;
            counter++;
        }
    }
    return results;

}

char *int2char(int iNumber) {
    int iNumbersCount = 0;
    int iTmpNum = iNumber;
    while (iTmpNum) {
        iTmpNum /= 10;
        iNumbersCount++;
    }
    char *buffer = malloc(iNumbersCount + 1);
    for (int i = iNumbersCount - 1; i >= 0; i--) {
        buffer[i] = (char) ((iNumber % 10) | 48);
        iNumber /= 10;
    }
    buffer[iNumbersCount] = '\0';
    return buffer;

}

/*int __put_state_wrapper(const char *key, const char *value) {
    int keySize;
    for (keySize = 0; key[keySize] != '\0'; ++keySize);
    int valueSize;
    for (valueSize = 0; value[valueSize] != '\0'; ++valueSize);
    return __put_state(key, keySize, value, valueSize);
}*/

/* convert character array to integer */
/*int char2int (char *array, int n){
    int number = 0;
    int mult = 1;

    n = (int)n < 0 ? -n : n;       *//* quick absolute value check  *//*

    *//* for each character in array *//*
    while (n--)
    {
        *//* if not digit or '-', check if number > 0, break or continue *//*
        if ((array[n] < '0' || array[n] > '9') && array[n] != '-') {
            if (number)
                break;
            else
                continue;
        }

        if (array[n] == '-') {      *//* if '-' if number, negate, break *//*
            if (number) {
                number = -number;
                break;
            }
        }
        else {                      *//* convert digit to numeric value   *//*
            number += (array[n] - '0') * mult;
            mult *= 10;
        }
    }

    return number;
}*/