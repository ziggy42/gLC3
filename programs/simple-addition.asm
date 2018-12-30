.ORIG x3000
    LD  R0, NUM                     ; load NUM in R0
    ADD R1, R0, R0                  ; NUM + NUM
    LD  R2, ASCII                   ; load the ascii offset into R2
    ADD R0, R1, R2                  ; Add ASCII offset to convert int to char
    OUT                             ; print the addition result

    LD R0, NEWLINE                  ; load \n in R0
    OUT                             ; output \n
HALT

NUM     .fill x02                   ; the number to add
ASCII   .fill x30                   ; ASCII offset
NEWLINE .fill x0A                   ; \n

.END
