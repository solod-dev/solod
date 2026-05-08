#include "main.h"

// -- Implementation --

so_int main_Person_Sleep(void* self) {
    main_Person* p = self;
    p->Age += 1;
    return p->Age;
}

int main(void) {
    main_Person p = (main_Person){.Name = so_str("Alice"), .Age = 30};
    main_Person_Sleep(&p);
    so_println("%.*s %s %" PRIdINT " %s", p.Name.len, p.Name.ptr, "is now", p.Age, "years old.");
    p.Nums[0] = 42;
    so_println("%s %" PRIdINT, "1st lucky number is", p.Nums[0]);
}
