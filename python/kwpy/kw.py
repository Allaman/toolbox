import time


def to_first_day_in_iso_week(year, week):
    # January 4th is always in ISO week 1
    fourth_jan = time.mktime((year, 1, 4, 0, 0, 0, 0, 0, 0))

    # Day of the week for January 4th (0=Monday, 6=Sunday)
    fourth_jan_weekday = time.localtime(fourth_jan)[6]

    # Calculate the start of week 1
    start_of_week1 = fourth_jan - (fourth_jan_weekday * 24 * 3600)

    # Calculate the start of the desired week
    start_of_desired_week = start_of_week1 + ((week - 1) * 7 * 24 * 3600)

    year, month, day, _, _, _, _, _, _ = time.localtime(start_of_desired_week)

    return (year, month, day)


def iso_week_number(year, month, day):
    t = time.mktime((year, month, day, 0, 0, 0, 0, 0, 0))

    # Calculate the weekday (0 = Monday, 6 = Sunday) and day of the year
    weekday = time.localtime(t)[6]

    # If the given date is January 1st
    if month == 1 and day == 1:
        # If the weekday is Sunday, it belongs to the last week of the previous year
        if weekday == 6:
            prev_year = year - 1
            # If previous year was a leap year and the weekday of its December 31st was Thursday, then week is 53
            if (prev_year % 4 == 0 and prev_year % 100 != 0) or (prev_year % 400 == 0):
                if (
                    time.localtime(time.mktime((prev_year, 12, 31, 0, 0, 0, 0, 0, 0)))[
                        6
                    ]
                    == 3
                ):
                    return 53
            return 52
        else:
            return 1

    # For the end of December
    if month == 12 and day >= 29:
        # If December 29th or 30th is a Monday, or if December 31st is a Monday in a leap year, then it's week 1 of next year
        if (
            (day == 29 and weekday == 0)
            or (day == 30 and weekday == 0)
            or (
                day == 31
                and weekday == 0
                and ((year % 4 == 0 and year % 100 != 0) or (year % 400 == 0))
            )
        ):
            return 1
        # If December 31st is a weekday other than Sunday or Monday, it belongs to the last week of the same year
        elif day == 31 and weekday == 6:  # Sunday
            return 52
        elif day == 31:
            if weekday == 3:  # Thursday
                return 53
            else:
                return 52

    # General case
    start_of_year = time.mktime((year, 1, 1, 0, 0, 0, 0, 0, 0))
    return int(((t - start_of_year) // (86400 * 7)) + 1)
