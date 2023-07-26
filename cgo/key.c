#include "devmem.h"
#include "key.h"
//      pd10   pd11  pd12   pd13
//pd5
//pd6
//pd7
//pd8
//pd9

void key_init(){
    Writebit(0x6C, 20, 0); //pd5 in
    Writebit(0x6C, 21, 0);
    Writebit(0x6C, 22, 0);
    Writebit(0x6C, 24, 0); //pd6 in
    Writebit(0x6C, 25, 0);
    Writebit(0x6C, 26, 0);
    Writebit(0x6C, 28, 0); //pd7 in
    Writebit(0x6C, 29, 0);
    Writebit(0x6C, 30, 0);
    Writebit(0x70, 0, 0); //pd8 in
    Writebit(0x70, 1, 0);
    Writebit(0x70, 2, 0);
    Writebit(0x70, 4, 0); //pd9 in
    Writebit(0x70, 5, 0);
    Writebit(0x70, 6, 0);


    Writebit(0x70, 8, 1); //pd10 out
    Writebit(0x70, 9, 0);
    Writebit(0x70, 10, 0);
    Writebit(0x70, 12, 1); //pd11 out
    Writebit(0x70, 13, 0);
    Writebit(0x70, 14, 0);

    Writebit(0x70, 16, 1); //pd12
    Writebit(0x70, 17, 0);
    Writebit(0x70, 18, 0);
    Writebit(0x70, 20, 1); //pd13
    Writebit(0x70, 21, 0);
    Writebit(0x70, 22, 0);

    Writebit(0x7C, 5, 1);
    Writebit(0x7C, 6, 1);
    Writebit(0x7C, 7, 1);
    Writebit(0x7C, 8, 1);
    Writebit(0x7C, 9, 1);

    Writebit(0x7C, 10, 1);
    Writebit(0x7C, 11, 1);
    Writebit(0x7C, 12, 1);  
    Writebit(0x7C, 13, 1);  

}

void set_line(unsigned char value){
    Writebit(0x7C, 5, 1);
    Writebit(0x7C, 6, 1);
    Writebit(0x7C, 7, 1);
    Writebit(0x7C, 8, 1);
    Writebit(0x7C, 9, 1);

    Writebit(0x7C, 10, 1);
    Writebit(0x7C, 11, 1);
    Writebit(0x7C, 12, 1);  
    Writebit(0x7C, 13, 1); 


   //逐行设置低电平
    switch(value){
        case 1:Writebit(0x7C, 10, 0);break;
        case 2:Writebit(0x7C, 11, 0);break;
        case 3:Writebit(0x7C, 12, 0);break;
        case 4:Writebit(0x7C, 13, 0);break;
        case 0:                         //全低
            Writebit(0x7C, 10, 0);
            Writebit(0x7C, 11, 0);
            Writebit(0x7C, 12, 0);  
            Writebit(0x7C, 13, 0); 
        break;
    }
}

unsigned int  get_line(){
    unsigned int value=0x00F;
    value|=Readbit(0x7C, 5)<<8;
    value|=Readbit(0x7C, 6)<<7;
    value|=Readbit(0x7C, 7)<<6;
    value|=Readbit(0x7C, 8)<<5;
    value|=Readbit(0x7C, 9)<<4;
    return value;
}

unsigned int read_key(){
    static unsigned int key_state=0,key_value,key_line;
    unsigned int key_result=NO_KEY,i;

    switch(key_state){
        case 0:
        key_line=0b000000001;
            for(i=1;i<=4;i++){
                set_line(i);
                key_value=KEY_MASK&get_line();
                if(key_value==KEY_MASK){  //no key press
                    key_line<<=1;
                }else{
                    key_state++;
                    break;
                }
            }
            break;
        case 1:
            if(key_value==(KEY_MASK&get_line())){
                 key_result=key_value|key_line;
//                switch(key_line|key_value){
//                    case 0b11100001:
//                    key_result=0b11100001;break;
//                }
                key_state++;
            }else{
                key_state--;
            }
        case 2:
            set_line(0);
            if((KEY_MASK&get_line())==KEY_MASK){
                key_state=0;
            }
            break;
    }
    return key_result;
}
