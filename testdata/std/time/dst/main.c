#include "main.h"

// -- Forward declarations --
static void format(void);
static void parse(void);
static void times(void);

// -- format.go --

static void format(void) {
    so_Slice buf = so_make_slice(so_byte, 64, 64);
    time_Time t = time_Date(2024, time_March, 15, 14, 30, 45, 0, time_UTC);
    {
        // RFC3339.
        so_String s = time_Time_Format(t, buf, time_RFC3339, time_UTC);
        if (so_string_ne(s, so_str("2024-03-15T14:30:45Z"))) {
            so_panic("unexpected RFC3339 format");
        }
    }
    {
        // RFC3339Nano.
        t = time_Date(2024, time_March, 15, 14, 30, 45, 123456789, time_UTC);
        so_String s = time_Time_Format(t, buf, time_RFC3339Nano, time_UTC);
        if (so_string_ne(s, so_str("2024-03-15T14:30:45.123456789Z"))) {
            so_panic("unexpected RFC3339Nano format");
        }
    }
    {
        // DateTime.
        so_String s = time_Time_Format(t, buf, time_DateTime, time_UTC);
        if (so_string_ne(s, so_str("2024-03-15 14:30:45"))) {
            so_panic("unexpected DateTime format");
        }
    }
    {
        // DateOnly.
        so_String s = time_Time_Format(t, buf, time_DateOnly, time_UTC);
        if (so_string_ne(s, so_str("2024-03-15"))) {
            so_panic("unexpected DateOnly format");
        }
    }
    {
        // TimeOnly.
        so_String s = time_Time_Format(t, buf, time_TimeOnly, time_UTC);
        if (so_string_ne(s, so_str("14:30:45"))) {
            so_panic("unexpected TimeOnly format");
        }
    }
    {
        // Custom format.
        so_String s = time_Time_Format(t, buf, so_str("%d.%m.%Y"), time_UTC);
        if (so_string_ne(s, so_str("15.03.2024"))) {
            so_panic("unexpected custom format");
        }
    }
    {
        // Time.String.
        so_String s = time_Time_String(t, buf);
        if (so_string_ne(s, so_str("2024-03-15T14:30:45Z"))) {
            so_panic("unexpected String format");
        }
    }
}

// -- main.go --

int main(void) {
    times();
    format();
    parse();
}

// -- parse.go --

static void parse(void) {
    // All tests use variants of 2024-03-15T14:30:45Z as the input time.
    {
        // RFC3339.
        time_TimeResult _res1 = time_Parse(time_RFC3339, so_str("2024-03-15T14:30:45Z"), 0);
        time_Time t = _res1.val;
        so_Error err = _res1.err;
        if (err != NULL) {
            so_panic("unexpected Parse RFC3339 error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 2024 || date.Month != time_March || date.Day != 15) {
            so_panic("unexpected Parse RFC3339 date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45) {
            so_panic("unexpected Parse RFC3339 clock");
        }
    }
    {
        // RFC3339Nano.
        time_TimeResult _res2 = time_Parse(time_RFC3339Nano, so_str("2024-03-15T14:30:45.123456789Z"), 0);
        time_Time t = _res2.val;
        so_Error err = _res2.err;
        if (err != NULL) {
            so_panic("unexpected Parse RFC3339Nano error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 2024 || date.Month != time_March || date.Day != 15) {
            so_panic("unexpected Parse RFC3339Nano date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45) {
            so_panic("unexpected Parse RFC3339Nano clock");
        }
        if (time_Time_Nanosecond(t) != 123456789) {
            so_panic("unexpected Parse RFC3339Nano nanosecond");
        }
    }
    {
        // RFC3339 with positive offset.
        // 14:30:45+05:00 is 09:30:45 UTC.
        time_TimeResult _res3 = time_Parse(time_RFC3339, so_str("2024-03-15T14:30:45+05:00"), 0);
        time_Time t = _res3.val;
        so_Error err = _res3.err;
        if (err != NULL) {
            so_panic("unexpected Parse RFC3339+offset error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 2024 || date.Month != time_March || date.Day != 15) {
            so_panic("unexpected Parse RFC3339+offset date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 9 || clock.Minute != 30 || clock.Second != 45) {
            so_panic("unexpected Parse RFC3339+offset clock");
        }
    }
    {
        // RFC3339 with negative offset.
        // 14:30:45-03:00 is 17:30:45 UTC.
        time_TimeResult _res4 = time_Parse(time_RFC3339, so_str("2024-03-15T14:30:45-03:00"), 0);
        time_Time t = _res4.val;
        so_Error err = _res4.err;
        if (err != NULL) {
            so_panic("unexpected Parse RFC3339-offset error");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 17 || clock.Minute != 30 || clock.Second != 45) {
            so_panic("unexpected Parse RFC3339-offset clock");
        }
    }
    {
        // RFC3339Nano with offset.
        // 14:30:45+05:30 is 09:00:45 UTC.
        time_TimeResult _res5 = time_Parse(time_RFC3339Nano, so_str("2024-03-15T14:30:45.123456789+05:30"), 0);
        time_Time t = _res5.val;
        so_Error err = _res5.err;
        if (err != NULL) {
            so_panic("unexpected Parse RFC3339Nano+offset error");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 9 || clock.Minute != 0 || clock.Second != 45) {
            so_panic("unexpected Parse RFC3339Nano+offset clock");
        }
        if (time_Time_Nanosecond(t) != 123456789) {
            so_panic("unexpected Parse RFC3339Nano+offset nanosecond");
        }
    }
    {
        // DateTime.
        time_TimeResult _res6 = time_Parse(time_DateTime, so_str("2024-03-15 14:30:45"), time_UTC);
        time_Time t = _res6.val;
        so_Error err = _res6.err;
        if (err != NULL) {
            so_panic("unexpected Parse DateTime error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 2024 || date.Month != time_March || date.Day != 15) {
            so_panic("unexpected Parse DateTime date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45) {
            so_panic("unexpected Parse DateTime clock");
        }
    }
    {
        // DateTime with offset parameter.
        // 14:30:45+05:30 is 09:00:45 UTC.
        // UTC+5:30
        time_Offset offset = (time_Offset)(5 * 3600 + 30 * 60);
        time_TimeResult _res7 = time_Parse(time_DateTime, so_str("2024-03-15 14:30:45"), offset);
        time_Time t = _res7.val;
        so_Error err = _res7.err;
        if (err != NULL) {
            so_panic("unexpected Parse DateTime+offset error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 2024 || date.Month != time_March || date.Day != 15) {
            so_panic("unexpected Parse DateTime+offset date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 9 || clock.Minute != 0 || clock.Second != 45) {
            so_panic("unexpected Parse DateTime+offset clock");
        }
    }
    {
        // DateOnly.
        time_TimeResult _res8 = time_Parse(time_DateOnly, so_str("2024-03-15"), time_UTC);
        time_Time t = _res8.val;
        so_Error err = _res8.err;
        if (err != NULL) {
            so_panic("unexpected Parse DateOnly error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 2024 || date.Month != time_March || date.Day != 15) {
            so_panic("unexpected Parse DateOnly date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 0 || clock.Minute != 0 || clock.Second != 0) {
            so_panic("unexpected Parse DateOnly clock");
        }
    }
    {
        // TimeOnly.
        time_TimeResult _res9 = time_Parse(time_TimeOnly, so_str("14:30:45"), time_UTC);
        time_Time t = _res9.val;
        so_Error err = _res9.err;
        if (err != NULL) {
            so_panic("unexpected Parse TimeOnly error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 0 || date.Month != time_January || date.Day != 1) {
            so_panic("unexpected Parse TimeOnly date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45) {
            so_panic("unexpected Parse TimeOnly clock");
        }
    }
    {
        // Custom format.
        time_TimeResult _res10 = time_Parse(so_str("%d.%m.%Y"), so_str("15.03.2024"), time_UTC);
        time_Time t = _res10.val;
        so_Error err = _res10.err;
        if (err != NULL) {
            so_panic("unexpected Parse custom error");
        }
        time_CalDate date = time_Time_Date(t, time_UTC);
        if (date.Year != 2024 || date.Month != time_March || date.Day != 15) {
            so_panic("unexpected Parse custom date");
        }
        time_CalClock clock = time_Time_Clock(t, time_UTC);
        if (clock.Hour != 0 || clock.Minute != 0 || clock.Second != 0) {
            so_panic("unexpected Parse custom clock");
        }
    }
    {
        // time.Parse error.
        time_TimeResult _res11 = time_Parse(so_str("%Y-%m-%d"), so_str("not-a-date"), time_UTC);
        so_Error err = _res11.err;
        if (err == NULL) {
            so_panic("expected Parse error");
        }
    }
}

// -- time.go --

static void times(void) {
    so_Slice buf = so_make_slice(so_byte, 64, 64);
    {
        // time.Date and time.Time properties.
        time_Time t = time_Date(2021, time_May, 10, 12, 33, 44, 777888999, time_UTC);
        if (time_Time_Year(t) != 2021) {
            so_panic("unexpected Time.Year");
        }
        if (time_Time_Month(t) != time_May) {
            so_panic("unexpected Time.Month");
        }
        if (time_Time_Day(t) != 10) {
            so_panic("unexpected Time.Day");
        }
        if (time_Time_Hour(t) != 12) {
            so_panic("unexpected Time.Hour");
        }
        if (time_Time_Minute(t) != 33) {
            so_panic("unexpected Time.Minute");
        }
        if (time_Time_Second(t) != 44) {
            so_panic("unexpected Time.Second");
        }
        if (time_Time_Nanosecond(t) != 777888999) {
            so_panic("unexpected Time.Nanosecond");
        }
        so_println("%" PRIdINT " %" PRIdINT " %" PRIdINT " %" PRIdINT " %" PRIdINT " %" PRIdINT " %" PRIdINT, time_Time_Year(t), time_Time_Month(t), time_Time_Day(t), time_Time_Hour(t), time_Time_Minute(t), time_Time_Second(t), time_Time_Nanosecond(t));
    }
    {
        // Time.Now.
        time_Time t = time_Now();
        if (time_Time_IsZero(t)) {
            so_panic("unexpected Time.IsZero");
        }
        so_println("%s %.*s", "UTC:", time_Time_String(t, buf).len, time_Time_String(t, buf).ptr);
        time_Offset utc5 = (time_Offset)(5 * 3600);
        so_println("%s %.*s", "UTC+5:", time_Time_Format(t, buf, time_RFC3339Nano, utc5).len, time_Time_Format(t, buf, time_RFC3339Nano, utc5).ptr);
    }
}
