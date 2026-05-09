#include "main.h"

// -- Forward declarations --
static void sliceTest(void);
static so_int descInt(void* a, void* b);
static void sortTest(void);
static void minMaxTest(void);

// -- main.go --

int main(void) {
    minMaxTest();
    sliceTest();
    sortTest();
    return 0;
}

// -- slice.go --

static void sliceTest(void) {
    {
        // Make a slice.
        so_Slice s = slices_Make(so_int, ((mem_Allocator){0}), (3));
        so_at(so_int, s, 0) = 11;
        so_at(so_int, s, 1) = 22;
        so_at(so_int, s, 2) = 33;
        if (so_len(s) != 3 || so_cap(s) != 3) {
            so_panic("Make failed");
        }
        if (so_at(so_int, s, 0) != 11 || so_at(so_int, s, 1) != 22 || so_at(so_int, s, 2) != 33) {
            so_panic("Make failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Append within capacity.
        so_Slice s = slices_MakeCap(so_int, ((mem_Allocator){0}), (0), (8));
        s = slices_Append(so_int, ((mem_Allocator){0}), (s), (10), (20), (30));
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 10 || so_at(so_int, s, 1) != 20 || so_at(so_int, s, 2) != 30) {
            so_panic("Append failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Append that triggers growth.
        so_Slice s = slices_MakeCap(so_int, ((mem_Allocator){0}), (0), (2));
        s = slices_Append(so_int, ((mem_Allocator){0}), (s), (1), (2));
        s = slices_Append(so_int, ((mem_Allocator){0}), (s), (3), (4), (5));
        if (so_len(s) != 5 || so_at(so_int, s, 0) != 1 || so_at(so_int, s, 4) != 5) {
            so_panic("Append grow failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Append to nil slice.
        so_Slice s = {0};
        s = slices_Append(so_int, ((mem_Allocator){0}), (s), (10), (20), (30));
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 10 || so_at(so_int, s, 1) != 20 || so_at(so_int, s, 2) != 30) {
            so_panic("Append to nil failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Extend from another slice.
        so_Slice s = slices_MakeCap(so_int, ((mem_Allocator){0}), (0), (8));
        so_Slice other = (so_Slice){(so_int[3]){100, 200, 300}, 3, 3};
        s = slices_Extend(so_int, ((mem_Allocator){0}), (s), (other));
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 100 || so_at(so_int, s, 2) != 300) {
            so_panic("Extend failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Extend a nil slice.
        so_Slice s = {0};
        so_Slice other = (so_Slice){(so_int[3]){10, 20, 30}, 3, 3};
        s = slices_Extend(so_int, ((mem_Allocator){0}), (s), (other));
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 10 || so_at(so_int, s, 1) != 20 || so_at(so_int, s, 2) != 30) {
            so_panic("Extend to nil failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Clone a slice.
        so_Slice s1 = (so_Slice){(so_int[3]){11, 22, 33}, 3, 3};
        so_Slice s2 = slices_Clone(so_int, ((mem_Allocator){0}), (s1));
        so_at(so_int, s2, 0) = 99;
        if (so_at(so_int, s1, 0) != 11 || so_at(so_int, s2, 0) != 99) {
            so_panic("Clone failed");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s2));
    }
    {
        // Equal slices.
        so_Slice s1 = (so_Slice){(so_int[3]){1, 2, 3}, 3, 3};
        so_Slice s2 = (so_Slice){(so_int[3]){1, 2, 3}, 3, 3};
        so_Slice s3 = (so_Slice){(so_int[3]){1, 2, 4}, 3, 3};
        so_Slice s4 = (so_Slice){(so_int[2]){1, 2}, 2, 2};
        so_Slice s5 = (so_Slice){0};
        so_Slice s6 = (so_Slice){0};
        if (!slices_Equal(so_int, (s1), (s2))) {
            so_panic("want s1 == s2");
        }
        if (slices_Equal(so_int, (s1), (s3))) {
            so_panic("want s1 != s3");
        }
        if (slices_Equal(so_int, (s1), (s4))) {
            so_panic("want s1 != s4");
        }
        if (!slices_Equal(so_int, (s5), (s6))) {
            so_panic("want empty and nil slices equal");
        }
    }
    {
        // Equal string slices.
        so_Slice s1 = (so_Slice){(so_String[3]){so_str("a"), so_str("b"), so_str("c")}, 3, 3};
        so_Slice s2 = (so_Slice){(so_String[3]){so_str("a"), so_str("b"), so_str("c")}, 3, 3};
        so_Slice s3 = (so_Slice){(so_String[3]){so_str("a"), so_str("b"), so_str("d")}, 3, 3};
        if (!slices_Equal(so_String, (s1), (s2))) {
            so_panic("want s1 == s2");
        }
        if (slices_Equal(so_String, (s1), (s3))) {
            so_panic("want s1 != s3");
        }
    }
    {
        // Equal struct slices.
        typedef struct point {
            so_int x;
            so_int y;
        } point;
        so_Slice s1 = (so_Slice){(point[2]){(point){1, 2}, (point){3, 4}}, 2, 2};
        so_Slice s2 = (so_Slice){(point[2]){(point){1, 2}, (point){3, 4}}, 2, 2};
        so_Slice s3 = (so_Slice){(point[2]){(point){1, 2}, (point){3, 5}}, 2, 2};
        if (!slices_Equal(point, (s1), (s2))) {
            so_panic("want s1 == s2");
        }
        if (slices_Equal(point, (s1), (s3))) {
            so_panic("want s1 != s3");
        }
    }
    {
        // Index of an element.
        so_Slice ints = (so_Slice){(so_int[4]){10, 20, 30, 20}, 4, 4};
        if (slices_Index(so_int, (ints), (20)) != 1) {
            so_panic("Index failed");
        }
        if (slices_Index(so_int, (ints), (40)) != -1) {
            so_panic("Index failed");
        }
        so_Slice strs = (so_Slice){(so_String[4]){so_str("a"), so_str("b"), so_str("c"), so_str("b")}, 4, 4};
        if (slices_Index(so_String, (strs), (so_str("b"))) != 1) {
            so_panic("Index failed");
        }
        if (slices_Index(so_String, (strs), (so_str("d"))) != -1) {
            so_panic("Index failed");
        }
    }
    {
        // Contains an element.
        so_Slice ints = (so_Slice){(so_int[4]){10, 20, 30, 20}, 4, 4};
        if (!slices_Contains(so_int, (ints), (20))) {
            so_panic("Contains failed");
        }
        if (slices_Contains(so_int, (ints), (40))) {
            so_panic("Contains failed");
        }
        so_Slice strs = (so_Slice){(so_String[4]){so_str("a"), so_str("b"), so_str("c"), so_str("b")}, 4, 4};
        if (!slices_Contains(so_String, (strs), (so_str("b")))) {
            so_panic("Contains failed");
        }
        if (slices_Contains(so_String, (strs), (so_str("d")))) {
            so_panic("Contains failed");
        }
    }
}

// -- sort.go --

static so_int descInt(void* a, void* b) {
    so_int va = *(so_int*)a;
    so_int vb = *(so_int*)b;
    return vb - va;
}

static void sortTest(void) {
    so_int ints[13] = {74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586};
    double float64s[15] = {74.3, 59.0, 238.2, -784.0, 2.3, 9845.768, -959.7485, 905, 7.8, 7.8, 74.3, 59.0, 238.2, -784.0, 2.3};
    so_String strs[8] = {so_str(""), so_str("Hello"), so_str("foo"), so_str("bar"), so_str("foo"), so_str("f00"), so_str("%*&^*&^&"), so_str("***")};
    {
        // IsSorted: false on unsorted data.
        if (slices_IsSorted(so_int, (so_array_slice(so_int, ints, 0, 13, 13)))) {
            so_panic("IsSorted: unsorted ints");
        }
        if (slices_IsSorted(so_String, (so_array_slice(so_String, strs, 0, 8, 8)))) {
            so_panic("IsSorted: unsorted strs");
        }
        // IsSorted: true on sorted data.
        so_Slice sorted = (so_Slice){(so_int[5]){1, 2, 3, 4, 5}, 5, 5};
        if (!slices_IsSorted(so_int, (sorted))) {
            so_panic("IsSorted: sorted ints");
        }
        so_Slice sortedStrs = (so_Slice){(so_String[3]){so_str("a"), so_str("b"), so_str("c")}, 3, 3};
        if (!slices_IsSorted(so_String, (sortedStrs))) {
            so_panic("IsSorted: sorted strs");
        }
    }
    {
        // IsSortedFunc: false on unsorted data.
        cmp_Func compare = cmp_FuncFor(so_int);
        if (slices_IsSortedFunc(so_int, (so_array_slice(so_int, ints, 0, 13, 13)), (compare))) {
            so_panic("IsSortedFunc: unsorted ints");
        }
        // IsSortedFunc: true on sorted data.
        so_Slice sorted = (so_Slice){(so_int[5]){1, 2, 3, 4, 5}, 5, 5};
        if (!slices_IsSortedFunc(so_int, (sorted), (compare))) {
            so_panic("IsSortedFunc: sorted ints");
        }
    }
    {
        // Sort ints.
        so_Slice s = slices_Clone(so_int, ((mem_Allocator){0}), (so_array_slice(so_int, ints, 0, 13, 13)));
        slices_Sort(so_int, (s));
        if (!slices_IsSorted(so_int, (s))) {
            so_panic("Sort ints: not sorted");
        }
        if (so_at(so_int, s, 0) != -5467984 || so_at(so_int, s, 12) != 9845) {
            so_panic("Sort ints: wrong values");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // Sort float64s.
        so_Slice s = slices_Clone(double, ((mem_Allocator){0}), (so_array_slice(double, float64s, 0, 15, 15)));
        slices_Sort(double, (s));
        if (!slices_IsSorted(double, (s))) {
            so_panic("Sort float64s: not sorted");
        }
        if (so_at(double, s, 0) != -959.7485 || so_at(double, s, 14) != 9845.768) {
            so_panic("Sort float64s: wrong values");
        }
        slices_Free(double, ((mem_Allocator){0}), (s));
    }
    {
        // Sort strings.
        so_Slice s = slices_Clone(so_String, ((mem_Allocator){0}), (so_array_slice(so_String, strs, 0, 8, 8)));
        slices_Sort(so_String, (s));
        if (!slices_IsSorted(so_String, (s))) {
            so_panic("Sort strings: not sorted");
        }
        if (so_string_ne(so_at(so_String, s, 0), so_str("")) || so_string_ne(so_at(so_String, s, 7), so_str("foo"))) {
            so_panic("Sort strings: wrong values");
        }
        slices_Free(so_String, ((mem_Allocator){0}), (s));
    }
    {
        // SortFunc (reverse order).
        so_Slice s = slices_Clone(so_int, ((mem_Allocator){0}), (so_array_slice(so_int, ints, 0, 13, 13)));
        slices_SortFunc(so_int, (s), (descInt));
        if (!slices_IsSortedFunc(so_int, (s), (descInt))) {
            so_panic("SortFunc ints: not sorted");
        }
        if (so_at(so_int, s, 0) != 9845 || so_at(so_int, s, 12) != -5467984) {
            so_panic("SortFunc ints: wrong values");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // SortFunc with nil compare.
        typedef struct point {
            so_int x;
            so_int y;
        } point;
        so_Slice s = (so_Slice){(point[3]){(point){1, 2}, (point){3, 4}, (point){2, 3}}, 3, 3};
        slices_SortFunc(point, (s), (NULL));
        if (!slices_IsSortedFunc(point, (s), (NULL))) {
            so_panic("SortFunc with nil: not sorted");
        }
        if (so_at(point, s, 0).x != 1 || so_at(point, s, 0).y != 2) {
            so_panic("SortFunc with nil: wrong s[0]");
        }
        if (so_at(point, s, 1).x != 2 || so_at(point, s, 1).y != 3) {
            so_panic("SortFunc with nil: wrong s[1]");
        }
        if (so_at(point, s, 2).x != 3 || so_at(point, s, 2).y != 4) {
            so_panic("SortFunc with nil: wrong s[2]");
        }
    }
    {
        // SortStableFunc ints.
        so_Slice s = slices_Clone(so_int, ((mem_Allocator){0}), (so_array_slice(so_int, ints, 0, 13, 13)));
        cmp_Func compare = cmp_FuncFor(so_int);
        slices_SortStableFunc(so_int, (s), (compare));
        if (!slices_IsSorted(so_int, (s))) {
            so_panic("SortStable ints: not sorted");
        }
        if (so_at(so_int, s, 0) != -5467984 || so_at(so_int, s, 12) != 9845) {
            so_panic("SortStable ints: wrong values");
        }
        slices_Free(so_int, ((mem_Allocator){0}), (s));
    }
    {
        // SortStableFunc float64s.
        so_Slice s = slices_Clone(double, ((mem_Allocator){0}), (so_array_slice(double, float64s, 0, 15, 15)));
        cmp_Func compare = cmp_FuncFor(double);
        slices_SortStableFunc(double, (s), (compare));
        if (!slices_IsSorted(double, (s))) {
            so_panic("SortStable float64s: not sorted");
        }
        if (so_at(double, s, 0) != -959.7485 || so_at(double, s, 14) != 9845.768) {
            so_panic("SortStable float64s: wrong values");
        }
        slices_Free(double, ((mem_Allocator){0}), (s));
    }
    {
        // SortStableFunc strings.
        so_Slice s = slices_Clone(so_String, ((mem_Allocator){0}), (so_array_slice(so_String, strs, 0, 8, 8)));
        cmp_Func compare = cmp_FuncFor(so_String);
        slices_SortStableFunc(so_String, (s), (compare));
        if (!slices_IsSorted(so_String, (s))) {
            so_panic("SortStable strings: not sorted");
        }
        if (so_string_ne(so_at(so_String, s, 0), so_str("")) || so_string_ne(so_at(so_String, s, 7), so_str("foo"))) {
            so_panic("SortStable strings: wrong values");
        }
        slices_Free(so_String, ((mem_Allocator){0}), (s));
    }
}

static void minMaxTest(void) {
    {
        // Min and Max on ints.
        so_Slice ints = (so_Slice){(so_int[6]){3, 1, 4, 1, 5, 9}, 6, 6};
        if (slices_Min(so_int, (ints)) != 1) {
            so_panic("Min ints: wrong value");
        }
        if (slices_Max(so_int, (ints)) != 9) {
            so_panic("Max ints: wrong value");
        }
    }
    {
        // Min and Max on strings.
        so_Slice strs = (so_Slice){(so_String[3]){so_str("banana"), so_str("apple"), so_str("cherry")}, 3, 3};
        if (so_string_ne(slices_Min(so_String, (strs)), so_str("apple"))) {
            so_panic("Min strings: wrong value");
        }
        if (so_string_ne(slices_Max(so_String, (strs)), so_str("cherry"))) {
            so_panic("Max strings: wrong value");
        }
    }
}
