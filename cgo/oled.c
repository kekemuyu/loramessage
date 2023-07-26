#include "font.h"
#include "devmem.h"
#include <linux/unistd.h>
#include "oled.h"
//OLED的显存
//存攠�式如下.
//[0]0 1 2 3 ... 127
//[1]0 1 2 3 ... 127
//[2]0 1 2 3 ... 127
//[3]0 1 2 3 ... 127
//[4]0 1 2 3 ... 127
//[5]0 1 2 3 ... 127
//[6]0 1 2 3 ... 127
//[7]0 1 2 3 ... 127

u8 OLED_GRAM[144][8];

void OLED_SCLK_Clr() //CLK
{
  //Writebit(0xa0, 2, 0)  ;
    Writebit(0x7C, 3, 0)  ;
}
void OLED_SCLK_Set()
{
//   Writebit(0xa0, 2, 1) ;
 Writebit(0x7C, 3, 1)  ;
}
void OLED_SDIN_Clr() //DIN
{
//    Writebit(0xa0, 3, 0);
 Writebit(0x7C, 4, 0)  ;
}
void OLED_SDIN_Set()
{
//    Writebit(0xa0, 3, 1);
 Writebit(0x7C, 4, 1)  ;
}


void OLED_RST_Clr() //GPIO_ResetBits(GPIOD,GPIO_Pin_4)//RES
{
//    Writebit(0xa0, 4, 0);
 Writebit(0x7C, 1, 0)  ;
}
void OLED_RST_Set() //GPIO_SetBits(GPIOD,GPIO_Pin_4)
{
//    Writebit(0xa0, 4, 1);
 Writebit(0x7C, 1, 1)  ;
}
void OLED_DC_Clr()  //GPIO_ResetBits(GPIOD,GPIO_Pin_5)//DC
{
//  Writebit(0xa0, 5, 0)  ;
 Writebit(0x7C, 2, 0)  ;
}
void OLED_DC_Set()  //GPIO_SetBits(GPIOD,GPIO_Pin_5)
{
//    Writebit(0xa0, 5, 1);
 Writebit(0x7C, 2, 1)  ;
}

void OLED_CS_Clr()  //CS
{
//    Writebit(0xa0, 6, 0);
 Writebit(0x7C, 0, 0)  ;
}
void OLED_CS_Set()  //GPIO_SetBits(GPIOD,GPIO_Pin_3)
{
//    Writebit(0xa0, 6, 1);
 Writebit(0x7C, 0, 1)  ;
}
#if OLED_MODE==1

//向SSD1106写入一个字节。
//dat:要写入的数据/命令
//cmd:数据/命令标志 0,表示命令;1,表示数据;
void OLED_WR_Byte(u8 dat,u8 cmd)
{
    DATAOUT(dat);
    if(cmd)
      OLED_DC_Set();
    else
      OLED_DC_Clr();
    OLED_CS_Clr();
    OLED_WR_Clr();
    OLED_WR_Set();
    OLED_CS_Set();
    OLED_DC_Set();
}
#else
//向SSD1106写入一个字节。
//dat:要写入的数据/命令
//cmd:数据/命令标志 0,表示命令;1,表示数据;
void OLED_WR_Byte(u8 dat,u8 cmd)
{
    u8 i;
    if(cmd)
      OLED_DC_Set();
    else
      OLED_DC_Clr();
    OLED_CS_Clr();
    for(i=0;i<8;i++)
    {
        OLED_SCLK_Clr();
        if(dat&0x80)
           OLED_SDIN_Set();
        else
           OLED_SDIN_Clr();
        OLED_SCLK_Set();
        dat<<=1;
    }
    OLED_CS_Set();
    OLED_DC_Set();
}
#endif
void OLED_Set_Pos(unsigned char x, unsigned char y)
{
    OLED_WR_Byte(0xb0+y,OLED_CMD);
    OLED_WR_Byte(((x&0xf0)>>4)|0x10,OLED_CMD);
    OLED_WR_Byte((x&0x0f)|0x01,OLED_CMD);
}
//开启OLED显示
void OLED_Display_On(void)
{
    OLED_WR_Byte(0X8D,OLED_CMD);  //SET DCDC命令
    OLED_WR_Byte(0X14,OLED_CMD);  //DCDC ON
    OLED_WR_Byte(0XAF,OLED_CMD);  //DISPLAY ON
}
//关闭OLED显示
void OLED_Display_Off(void)
{
    OLED_WR_Byte(0X8D,OLED_CMD);  //SET DCDC命令
    OLED_WR_Byte(0X10,OLED_CMD);  //DCDC OFF
    OLED_WR_Byte(0XAE,OLED_CMD);  //DISPLAY OFF
}


//清屏函数,清完屏,整个屏幕是黑色的!和没点亮一样!!!
void OLED_Clear(void)
{
	u8 i,n;
	for(i=0;i<8;i++)
	{
	   for(n=0;n<128;n++)
			{
			 OLED_GRAM[n][i]=0;//Çå³ýËùÓÐÊý¾Ý
			}
  }
}

void OLED_Unclear(void)
{
	u8 i,n;
	for(i=0;i<8;i++)
	{
	   for(n=0;n<128;n++)
			{
			 OLED_GRAM[n][i]=1;//Çå³ýËùÓÐÊý¾Ý
			}
  }
}

void OLED_Show(){
    u8 i,n;
    for(i=0;i<8;i++)
    {
       OLED_WR_Byte(0xb0+i,OLED_CMD); 
       OLED_WR_Byte(0x00,OLED_CMD);   
       OLED_WR_Byte(0x10,OLED_CMD);   
       for(n=0;n<128;n++)
    	 OLED_WR_Byte(OLED_GRAM[n][i],OLED_DATA);
    }
}
//画点 
//x:0~127
//y:0~63
//t:1 填充 0 清空	
void OLED_Point(u8 x,u8 y,u8 col){
   	u8 i,m,n;
	i=y/8;
	m=y%8;
	n=1<<m;
	if(col){OLED_GRAM[x][i]|=n;}
	else
	{
		OLED_GRAM[x][i]=~OLED_GRAM[x][i];
		OLED_GRAM[x][i]|=n;
		OLED_GRAM[x][i]=~OLED_GRAM[x][i];
	} 
}

void OLED_Line(u8 x,u8 y,u8 x2,u8 y2,u8 col){
    
    
}


void OLED_FillRect(u8 x,u8 y,u8 w,u8 h,u8 col){
    for(u8 i=0;i<h;i++){
        for(u8 j=0;j<w;j++){
           OLED_Point(x+j,y+i,col);  
        }
    } 
}

void OLED_Rect(u8 x,u8 y,u8 w,u8 h,u8 col){
    
    
}
//在指定位置显示一个字符,包括部分字符
//x:0~127
//y:0~63
//mode:0,反白显示;1,正常显示
//size:选择字体 16/12
void OLED_ShowChar(u8 x,u8 y,u8 chr,u8 col)
{
	u8 i,m,temp,chr1,size1=16;
	u8 x0=x,y0=y;

	chr1=chr-' ';  
	for(i=0;i<size1;i++)
	{
		temp=asc2_1608[chr1][i];
		for(m=0;m<8;m++)
		{
			if(temp&0x01)OLED_Point(x,y,col);
			else OLED_Point(x,y,!col);
			temp>>=1;
			y++;
		}
		x++;
		if((size1!=8)&&((x-x0)==size1/2))
		{x=x0;y0=y0+8;}
		y=y0;
  }
}

//m^n函数
u32 oled_pow(u8 m,u8 n)
{
    u32 result=1;
    while(n--)result*=m;
    return result;
}
//显示2个数字
//x,y :起点坐标
//len :数字的位数
//size:字体大小
//mode:模式    0,填充模式;1,叠加模式
//num:数值(0~4294967295);
void OLED_ShowNum(u8 x,u8 y,u32 num,u8 len,u8 size)
{
    u8 t,temp;
    u8 enshow=0;
    for(t=0;t<len;t++)
    {
        temp=(num/oled_pow(10,len-t-1))%10;
        if(enshow==0&&t<(len-1))
        {
            if(temp==0)
            {
                OLED_ShowChar(x+(size/2)*t,y,' ',1);
                continue;
            }else enshow=1;

        }
         OLED_ShowChar(x+(size/2)*t,y,temp+'0',1);
    }
}
//显示一个字符号串
void OLED_ShowString(u8 x,u8 y,char *chr,u8 col)
{
    u8 size1=16;
	while((*chr>=' ')&&(*chr<='~'))
	{
		OLED_ShowChar(x,y,*chr,col);
		if(size1==8)x+=6;
		else x+=size1/2;
		chr++;
  }
}

//显示汉字
void OLED_ShowBuffer(u8 x,u8 y,u8 *chr,u8 col)
{
    u8 m,temp;
    u8 x0=x,y0=y;
    u8 i;
    for(i=0;i<32;i++){
        temp=chr[i];
  
       for(m=0;m<8;m++)
		{
			if((temp&0x80)==0x80)OLED_Point(x,y,col);
			else OLED_Point(x,y,!col);
			temp<<=1;
			x++;
		} 
        if(i%2==0){
            x=x0+8;
 
        }else{
            x=x0;
           y++; 
        }
    }
}
/***********功能描述：显示显示BMP图片128×64起始点坐标(x,y),x的范围0～127，y为页的范围0～7*****************/
void OLED_DrawBMP(unsigned char x0, unsigned char y0,unsigned char x1, unsigned char y1, unsigned char *BMP)
{
 unsigned int j=0;
 unsigned char x,y;

  if(y1%8==0) y=y1/8;
  else y=y1/8+1;
    for(y=y0;y<y1;y++)
    {
        OLED_Set_Pos(x0,y);
    for(x=x0;x<x1;x++)
        {
            OLED_WR_Byte(BMP[j++],OLED_DATA);
        }
    }
}


//初始化SSD1306
void OLED_Init(void)
{
    Openfile(0x01C0C000);   //close tcon
    Writebit(0, 31, 0);
    Closefile();
    Openfile(GPIO_BASE);
//        Writebit(0x90, 8, 1); //pe2 out
//        Writebit(0x90, 9, 0);
//        Writebit(0x90, 10, 0);

//        Writebit(0x90, 12, 1); //pe3 out
//        Writebit(0x90, 13, 0);
//        Writebit(0x90, 14, 0);

//        Writebit(0x90, 16, 1); //pe4 out
//        Writebit(0x90, 17, 0);
//        Writebit(0x90, 18, 0);

//        Writebit(0x90, 20, 1); //pe5 out
//        Writebit(0x90, 21, 0);
//        Writebit(0x90, 22, 0);

//        Writebit(0x90, 24, 1) ;//pe6 out
//        Writebit(0x90, 25, 0);
//        Writebit(0x90, 26, 0) ;

       Writebit(0x6C, 0, 1); //pd0 out
       Writebit(0x6C, 1, 0);
       Writebit(0x6C, 2, 0);

       Writebit(0x6C, 4, 1); //pd1 out
       Writebit(0x6C, 5, 0);
       Writebit(0x6C, 6, 0);

       Writebit(0x6C, 8, 1); //pd2 out
       Writebit(0x6C, 9, 0);
       Writebit(0x6C, 10, 0);

       Writebit(0x6C, 12, 1); //pd3 out
       Writebit(0x6C, 13, 0);
       Writebit(0x6C, 14, 0);

       Writebit(0x6C, 16, 1) ;//pd4 out
       Writebit(0x6C, 17, 0);
       Writebit(0x6C, 18, 0) ;
    OLED_RST_Set();
    usleep(100000);
    OLED_RST_Clr();
    usleep(100000);
    OLED_RST_Set();

    OLED_WR_Byte(0xAE,OLED_CMD);//--turn off oled panel
    OLED_WR_Byte(0x00,OLED_CMD);//---set low column address
    OLED_WR_Byte(0x10,OLED_CMD);//---set high column address
    OLED_WR_Byte(0x40,OLED_CMD);//--set start line address  Set Mapping RAM Display Start Line (0x00~0x3F)
    OLED_WR_Byte(0x81,OLED_CMD);//--set contrast control register
    OLED_WR_Byte(0x66,OLED_CMD); // Set SEG Output Current Brightness
    OLED_WR_Byte(0xA1,OLED_CMD);//--Set SEG/Column Mapping     0xa0左右反置 0xa1正常
    OLED_WR_Byte(0xC8,OLED_CMD);//Set COM/Row Scan Direction   0xc0上下反置 0xc8正常
    OLED_WR_Byte(0xA6,OLED_CMD);//--set normal display
    OLED_WR_Byte(0xA8,OLED_CMD);//--set multiplex ratio(1 to 64)
    OLED_WR_Byte(0x3f,OLED_CMD);//--1/64 duty
    OLED_WR_Byte(0xD3,OLED_CMD);//-set display offset    Shift Mapping RAM Counter (0x00~0x3F)
    OLED_WR_Byte(0x00,OLED_CMD);//-not offset
    OLED_WR_Byte(0xd5,OLED_CMD);//--set display clock divide ratio/oscillator frequency
    OLED_WR_Byte(0x80,OLED_CMD);//--set divide ratio, Set Clock as 100 Frames/Sec
    OLED_WR_Byte(0xD9,OLED_CMD);//--set pre-charge period
    OLED_WR_Byte(0xF1,OLED_CMD);//Set Pre-Charge as 15 Clocks & Discharge as 1 Clock
    OLED_WR_Byte(0xDA,OLED_CMD);//--set com pins hardware configuration
    OLED_WR_Byte(0x12,OLED_CMD);
    OLED_WR_Byte(0xDB,OLED_CMD);//--set vcomh
    OLED_WR_Byte(0x40,OLED_CMD);//Set VCOM Deselect Level
    OLED_WR_Byte(0x20,OLED_CMD);//-Set Page Addressing Mode (0x00/0x01/0x02)
    OLED_WR_Byte(0x02,OLED_CMD);//
    OLED_WR_Byte(0x8D,OLED_CMD);//--set Charge Pump enable/disable
    OLED_WR_Byte(0x14,OLED_CMD);//--set(0x10) disable
    OLED_WR_Byte(0xA4,OLED_CMD);// Disable Entire Display On (0xa4/0xa5)
    OLED_WR_Byte(0xA6,OLED_CMD);// Disable Inverse Display On (0xa6/a7)
    OLED_WR_Byte(0xAF,OLED_CMD);//--turn on oled panel

    OLED_WR_Byte(0xAF,OLED_CMD); /*display ON*/
    OLED_Clear();
    OLED_Set_Pos(0,0);
}
