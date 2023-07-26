
#define u8 unsigned char
#define u32  unsigned long
#define OLED_MODE 0
#define SIZE 16
#define XLevelL        0x00
#define XLevelH        0x10
#define Max_Column    128
#define Max_Row        64
#define    Brightness    0xFF
#define X_WIDTH     128
#define Y_WIDTH     64


#define OLED_CMD  0     //写命令
#define OLED_DATA 1     //写数据
#define GPIO_BASE 0x1c20800
//OLED控制用函数
void OLED_CS_Clr();
void OLED_CS_Set();

void OLED_RST_Clr();
void OLED_RST_Set() ;

void OLED_DC_Clr()  ;
void OLED_DC_Set()  ;


void OLED_SCLK_Clr() ;
void OLED_SCLK_Set() ;

void OLED_SDIN_Clr() ;
void OLED_SDIN_Set() ;
void OLED_WR_Byte(u8 dat,u8 cmd);
void OLED_Display_On(void);
void OLED_Display_Off(void);                                                    
void OLED_Init(void);
void OLED_Clear(void);
void OLED_Unclear(void);
void OLED_DrawPoint(u8 x,u8 y,u8 t);
void OLED_Fill(u8 x1,u8 y1,u8 x2,u8 y2,u8 dot);
void OLED_ShowChar(u8 x,u8 y,u8 chr,u8 col);
void OLED_ShowNum(u8 x,u8 y,u32 num,u8 len,u8 size);
void OLED_ShowString(u8 x,u8 y,char *chr,u8 col);
void OLED_Set_Pos(unsigned char x, unsigned char y);
void OLED_ShowBuffer(u8 x,u8 y,u8 *chr,u8 col);
void OLED_DrawBMP(unsigned char x0, unsigned char y0,unsigned char x1, unsigned char y1, unsigned char *BMP);
void OLED_FillRect(u8 x,u8 y,u8 w,u8 h,u8 col);
void OLED_Point(u8 x,u8 y,u8 col);
void OLED_Show();