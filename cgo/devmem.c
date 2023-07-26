#include"devmem.h"


int tempfd;
void *map_base,*virt_addr;

void Openfile(long  target){
    int fd;

    if((tempfd = open("/dev/mem", O_RDWR | O_SYNC)) == -1) FATAL;
    fd = fcntl(tempfd, F_DUPFD, 0);
    if(fd<0){
            FATAL;
    }

    /* Map one page */
    map_base = mmap(0, MAP_SIZE, PROT_READ | PROT_WRITE, MAP_SHARED, fd, target & ~MAP_MASK);
    if(map_base == (void *) -1) FATAL;

    virt_addr = map_base + (target & MAP_MASK);
}

void Closefile(){
     if(munmap(map_base, MAP_SIZE) == -1) FATAL;
     close(tempfd);
}

void Writebit(int offset,int bitsize ,char value){
    unsigned long read_result, writeval;

    read_result = *((unsigned long *) (virt_addr+offset));


    if(value==0){
       read_result&=~(1<<bitsize);
    }else{
        read_result|=1<<bitsize;
    }

    writeval=read_result;
    *((unsigned long *) (virt_addr+offset))=writeval;
}


char Readbit(int offset,int bitsize){
    unsigned long read_result, writeval;

    read_result = *((unsigned long *) (virt_addr+offset));

    return (read_result>>bitsize)&0x0001;
}