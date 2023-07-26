#define NO_KEY 511
#define KEY_MASK 0b111110000

void key_init();
void set_line(unsigned char value);
unsigned int  get_line();
unsigned int read_key();
