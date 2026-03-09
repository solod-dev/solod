#include <stdint.h>
#include <time.h>

#define time_TimeT time_t
#define time_ClocksPerSec CLOCKS_PER_SEC

// Wrapper struct with So-style field names.
typedef struct {
    int Sec, Min, Hour, Mday, Mon, Year, Wday, Yday, Isdst;
} time_Tm;

static inline struct tm time_Tm_to_c(time_Tm t) {
    struct tm ct = {0};
    ct.tm_sec = t.Sec;
    ct.tm_min = t.Min;
    ct.tm_hour = t.Hour;
    ct.tm_mday = t.Mday;
    ct.tm_mon = t.Mon;
    ct.tm_year = t.Year;
    ct.tm_wday = t.Wday;
    ct.tm_yday = t.Yday;
    ct.tm_isdst = t.Isdst;
    return ct;
}

static inline time_Tm time_Tm_from_c(struct tm ct) {
    time_Tm t;
    t.Sec = ct.tm_sec;
    t.Min = ct.tm_min;
    t.Hour = ct.tm_hour;
    t.Mday = ct.tm_mday;
    t.Mon = ct.tm_mon;
    t.Year = ct.tm_year;
    t.Wday = ct.tm_wday;
    t.Yday = ct.tm_yday;
    t.Isdst = ct.tm_isdst;
    return t;
}

#define time_Time(timer) time(timer)
#define time_Clock() ((int)clock())
#define time_Difftime(end, start) difftime(end, start)

static inline time_Tm time_Gmtime(time_TimeT* timer) {
    time_Tm result = time_Tm_from_c(*gmtime(timer));
    return result;
}

static inline time_TimeT time_Mktime(time_Tm* t) {
    struct tm ct = time_Tm_to_c(*t);
    time_TimeT r = mktime(&ct);
    *t = time_Tm_from_c(ct);
    return r;
}

static inline int time_Strftime(uint8_t* buf, size_t maxsize, const char* format, time_Tm* t) {
    struct tm ct = time_Tm_to_c(*t);
    return (int)strftime((char*)buf, maxsize, format, &ct);
}
