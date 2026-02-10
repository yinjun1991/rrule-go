package rrule

import (
	"sort"
	"time"
)

type iterInfo struct {
	recurrence  *Recurrence
	lastyear    int
	lastmonth   time.Month
	yearlen     int
	nextyearlen int
	firstyday   time.Time
	yearweekday int
	mmask       []int
	mrange      []int
	mdaymask    []int
	nmdaymask   []int
	wdaymask    []int
	wnomask     []int
	nwdaymask   []int
	eastermask  []int
}

func (info *iterInfo) rebuild(year int, month time.Month) {
	// Every mask is 7 days longer to handle cross-year weekly periods.
	if year != info.lastyear {
		info.yearlen = 365 + isLeap(year)
		info.nextyearlen = 365 + isLeap(year+1)
		info.firstyday = time.Date(
			year, time.January, 1, 0, 0, 0, 0,
			info.recurrence.dtstart.Location())
		info.yearweekday = toPyWeekday(info.firstyday.Weekday())
		info.wdaymask = WDAYMASK[info.yearweekday:]
		if info.yearlen == 365 {
			info.mmask = M365MASK
			info.mdaymask = MDAY365MASK
			info.nmdaymask = NMDAY365MASK
			info.mrange = M365RANGE
		} else {
			info.mmask = M366MASK
			info.mdaymask = MDAY366MASK
			info.nmdaymask = NMDAY366MASK
			info.mrange = M366RANGE
		}
		if len(info.recurrence.byweekno) == 0 {
			info.wnomask = nil
		} else {
			info.wnomask = make([]int, info.yearlen+7)
			firstwkst := pymod(7-info.yearweekday+info.recurrence.wkst, 7)
			no1wkst := firstwkst
			var wyearlen int
			if no1wkst >= 4 {
				no1wkst = 0
				// Number of days in the year, plus the days we got from last year.
				wyearlen = info.yearlen + pymod(info.yearweekday-info.recurrence.wkst, 7)
			} else {
				// Number of days in the year, minus the days we left in last year.
				wyearlen = info.yearlen - no1wkst
			}
			div, mod := divmod(wyearlen, 7)
			numweeks := div + mod/4
			for _, n := range info.recurrence.byweekno {
				if n < 0 {
					n += numweeks + 1
				}
				if !(0 < n && n <= numweeks) {
					continue
				}
				var i int
				if n > 1 {
					i = no1wkst + (n-1)*7
					if no1wkst != firstwkst {
						i -= 7 - firstwkst
					}
				} else {
					i = no1wkst
				}
				for j := 0; j < 7; j++ {
					info.wnomask[i] = 1
					i++
					if info.wdaymask[i] == info.recurrence.wkst {
						break
					}
				}
			}
			if contains(info.recurrence.byweekno, 1) {
				// Check week number 1 of next year as well
				// TODO: Check -numweeks for next year.
				i := no1wkst + numweeks*7
				if no1wkst != firstwkst {
					i -= 7 - firstwkst
				}
				if i < info.yearlen {
					// If week starts in next year, we
					// don't care about it.
					for j := 0; j < 7; j++ {
						info.wnomask[i] = 1
						i++
						if info.wdaymask[i] == info.recurrence.wkst {
							break
						}
					}
				}
			}
			if no1wkst != 0 {
				// Check last week number of last year as
				// well. If no1wkst is 0, either the year
				// started on week start, or week number 1
				// got days from last year, so there are no
				// days from last year's last week number in
				// this year.
				var lnumweeks int
				if !contains(info.recurrence.byweekno, -1) {
					lyearweekday := toPyWeekday(time.Date(
						year-1, 1, 1, 0, 0, 0, 0,
						info.recurrence.dtstart.Location()).Weekday())
					lno1wkst := pymod(7-lyearweekday+info.recurrence.wkst, 7)
					lyearlen := 365 + isLeap(year-1)
					if lno1wkst >= 4 {
						lno1wkst = 0
						lnumweeks = 52 + pymod(lyearlen+pymod(lyearweekday-info.recurrence.wkst, 7), 7)/4
					} else {
						lnumweeks = 52 + pymod(info.yearlen-no1wkst, 7)/4
					}
				} else {
					lnumweeks = -1
				}
				if contains(info.recurrence.byweekno, lnumweeks) {
					for i := 0; i < no1wkst; i++ {
						info.wnomask[i] = 1
					}
				}
			}
		}
	}
	if len(info.recurrence.bynweekday) != 0 && (month != info.lastmonth || year != info.lastyear) {
		var ranges [][]int
		if info.recurrence.freq == YEARLY {
			if len(info.recurrence.bymonth) != 0 {
				for _, month := range info.recurrence.bymonth {
					ranges = append(ranges, info.mrange[month-1:month+1])
				}
			} else {
				ranges = [][]int{{0, info.yearlen}}
			}
		} else if info.recurrence.freq == MONTHLY {
			ranges = [][]int{info.mrange[month-1 : month+1]}
		}
		if len(ranges) != 0 {
			// Weekly frequency won't get here, so we may not
			// care about cross-year weekly periods.
			info.nwdaymask = make([]int, info.yearlen)
			for _, x := range ranges {
				first, last := x[0], x[1]
				last--
				for _, y := range info.recurrence.bynweekday {
					wday, n := y.weekday, y.n
					var i int
					if n < 0 {
						i = last + (n+1)*7
						i -= pymod(info.wdaymask[i]-wday, 7)
					} else {
						i = first + (n-1)*7
						i += pymod(7-info.wdaymask[i]+wday, 7)
					}
					if first <= i && i <= last {
						info.nwdaymask[i] = 1
					}
				}
			}
		}
	}
	if len(info.recurrence.byeaster) != 0 {
		info.eastermask = make([]int, info.yearlen+7)
		eyday := easter(year).YearDay() - 1
		for _, offset := range info.recurrence.byeaster {
			info.eastermask[eyday+offset] = 1
		}
	}
	info.lastyear = year
	info.lastmonth = month
}

func (info *iterInfo) calcDaySet(freq Frequency, year int, month time.Month, day int) (start, end int) {
	switch freq {
	case YEARLY:
		return 0, info.yearlen

	case MONTHLY:
		start, end = info.mrange[month-1], info.mrange[month]
		return start, end

	case WEEKLY:
		// We need to handle cross-year weeks here.
		i := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).YearDay() - 1
		start, end = i, i+1
		for j := 0; j < 7; j++ {
			i++
			// if (not (0 <= i < self.yearlen) or
			//     self.wdaymask[i] == self.rrule._wkst):
			//  This will cross the year boundary, if necessary.
			if info.wdaymask[i] == info.recurrence.wkst {
				break
			}

			end = i + 1
		}

		return start, end

	default:
		// DAILY, HOURLY, MINUTELY, SECONDLY:
		i := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).YearDay() - 1
		return i, i + 1
	}
}

func (info *iterInfo) fillTimeSet(set *[]time.Time, freq Frequency, hour, minute, second int) {
	switch freq {
	case HOURLY:
		prepareTimeSet(set, len(info.recurrence.byminute)*len(info.recurrence.bysecond))
		for _, minute := range info.recurrence.byminute {
			for _, second := range info.recurrence.bysecond {
				*set = append(*set, time.Date(1, 1, 1, hour, minute, second, 0, info.recurrence.dtstart.Location()))
			}
		}
		sort.Sort(timeSlice(*set))
	case MINUTELY:
		prepareTimeSet(set, len(info.recurrence.bysecond))
		for _, second := range info.recurrence.bysecond {
			*set = append(*set, time.Date(1, 1, 1, hour, minute, second, 0, info.recurrence.dtstart.Location()))
		}
		sort.Sort(timeSlice(*set))
	case SECONDLY:
		prepareTimeSet(set, 1)
		*set = append(*set, time.Date(1, 1, 1, hour, minute, second, 0, info.recurrence.dtstart.Location()))
	default:
		prepareTimeSet(set, 0)
	}
}

// rIterator is a iterator of RRule
type rIterator struct {
	year     int
	month    time.Month
	day      int
	hour     int
	minute   int
	second   int
	weekday  int
	ii       iterInfo
	timeset  []time.Time
	total    int
	count    int // A value of 0 means count is unlimited.
	remain   reusingRemainSlice
	finished bool
	dayset   []optInt
}

func (iterator *rIterator) generate() {
	if iterator.finished {
		return
	}

	r := iterator.ii.recurrence

	for iterator.remain.Len() == 0 {
		// Get dayset with the right frequency
		setStart, setEnd := iterator.ii.calcDaySet(r.freq, iterator.year, iterator.month, iterator.day)
		iterator.fillDaySetMonotonic(setStart, setEnd)

		dayset := iterator.dayset
		filtered := false

		// Do the "hard" work ;-)
		for dayIndex, day := range dayset {
			i := day.Int
			if len(r.bymonth) != 0 && !contains(r.bymonth, iterator.ii.mmask[i]) ||
				len(r.byweekno) != 0 && iterator.ii.wnomask[i] == 0 ||
				len(r.byweekday) != 0 && !contains(r.byweekday, iterator.ii.wdaymask[i]) ||
				len(iterator.ii.nwdaymask) != 0 && iterator.ii.nwdaymask[i] == 0 ||
				len(r.byeaster) != 0 && iterator.ii.eastermask[i] == 0 ||
				(len(r.bymonthday) != 0 || len(r.bynmonthday) != 0) &&
					!contains(r.bymonthday, iterator.ii.mdaymask[i]) &&
					!contains(r.bynmonthday, iterator.ii.nmdaymask[i]) ||
				len(r.byyearday) != 0 &&
					(i < iterator.ii.yearlen &&
						!contains(r.byyearday, i+1) &&
						!contains(r.byyearday, -iterator.ii.yearlen+i) ||
						i >= iterator.ii.yearlen &&
							!contains(r.byyearday, i+1-iterator.ii.yearlen) &&
							!contains(r.byyearday, -iterator.ii.nextyearlen+i-iterator.ii.yearlen)) {
				dayset[dayIndex].Defined = false
				filtered = true
			}
		}

		// Output results
		if len(r.bysetpos) != 0 && len(iterator.timeset) != 0 {
			var poslist []time.Time
			for _, pos := range r.bysetpos {
				var daypos, timepos int
				if pos < 0 {
					daypos, timepos = divmod(pos, len(iterator.timeset))
				} else {
					daypos, timepos = divmod(pos-1, len(iterator.timeset))
				}
				var temp []int
				for _, day := range dayset {
					if day.Defined {
						temp = append(temp, day.Int)
					}
				}
				i, err := pySubscript(temp, daypos)
				if err != nil {
					continue
				}
				timeTemp := iterator.timeset[timepos]
				dateYear, dateMonth, dateDay := iterator.ii.firstyday.AddDate(0, 0, i).Date()
				tempHour, tempMinute, tempSecond := timeTemp.Clock()
				res := time.Date(dateYear, dateMonth, dateDay,
					tempHour, tempMinute, tempSecond,
					timeTemp.Nanosecond(), timeTemp.Location())
				if !timeContains(poslist, res) {
					poslist = append(poslist, res)
				}
			}
			sort.Sort(timeSlice(poslist))
			for _, res := range poslist {
				if !r.until.IsZero() && res.After(r.until) {
					r.len = iterator.total
					iterator.finished = true
					return
				} else if !res.Before(r.dtstart) {
					iterator.total++
					iterator.remain.Append(res)
					if iterator.count > 0 {
						iterator.count--
						if iterator.count == 0 {
							r.len = iterator.total
							iterator.finished = true
							return
						}
					}
				}
			}
		} else {
			for _, day := range dayset {
				if !day.Defined {
					continue
				}
				i := day.Int
				dateYear, dateMonth, dateDay := iterator.ii.firstyday.AddDate(0, 0, i).Date()
				for _, timeTemp := range iterator.timeset {
					tempHour, tempMinute, tempSecond := timeTemp.Clock()
					res := time.Date(dateYear, dateMonth, dateDay,
						tempHour, tempMinute, tempSecond,
						timeTemp.Nanosecond(), timeTemp.Location())
					if !r.until.IsZero() && res.After(r.until) {
						r.len = iterator.total
						iterator.finished = true
						return
					} else if !res.Before(r.dtstart) {
						iterator.total++
						iterator.remain.Append(res)
						if iterator.count > 0 {
							iterator.count--
							if iterator.count == 0 {
								r.len = iterator.total
								iterator.finished = true
								return
							}
						}
					}
				}
			}
		}
		// Handle frequency and interval
		fixday := false
		if r.freq == YEARLY {
			iterator.year += r.interval
			if iterator.year > MAXYEAR {
				r.len = iterator.total
				iterator.finished = true
				return
			}
			iterator.ii.rebuild(iterator.year, iterator.month)
		} else if r.freq == MONTHLY {
			iterator.month += time.Month(r.interval)
			if iterator.month > 12 {
				div, mod := divmod(int(iterator.month), 12)
				iterator.month = time.Month(mod)
				iterator.year += div
				if iterator.month == 0 {
					iterator.month = 12
					iterator.year--
				}
				if iterator.year > MAXYEAR {
					r.len = iterator.total
					iterator.finished = true
					return
				}
			}
			iterator.ii.rebuild(iterator.year, iterator.month)
		} else if r.freq == WEEKLY {
			if r.wkst > iterator.weekday {
				iterator.day += -(iterator.weekday + 1 + (6 - r.wkst)) + r.interval*7
			} else {
				iterator.day += -(iterator.weekday - r.wkst) + r.interval*7
			}
			iterator.weekday = r.wkst
			fixday = true
		} else if r.freq == DAILY {
			iterator.day += r.interval
			fixday = true
		} else if r.freq == HOURLY {
			if filtered {
				// Jump to one iteration before next day
				iterator.hour += ((23 - iterator.hour) / r.interval) * r.interval
			}
			for {
				iterator.hour += r.interval
				div, mod := divmod(iterator.hour, 24)
				if div != 0 {
					iterator.hour = mod
					iterator.day += div
					fixday = true
				}
				if len(r.byhour) == 0 || contains(r.byhour, iterator.hour) {
					break
				}
			}
			iterator.ii.fillTimeSet(&iterator.timeset, r.freq, iterator.hour, iterator.minute, iterator.second)
		} else if r.freq == MINUTELY {
			if filtered {
				// Jump to one iteration before next day
				iterator.minute += ((1439 - (iterator.hour*60 + iterator.minute)) / r.interval) * r.interval
			}
			for {
				iterator.minute += r.interval
				div, mod := divmod(iterator.minute, 60)
				if div != 0 {
					iterator.minute = mod
					iterator.hour += div
					div, mod = divmod(iterator.hour, 24)
					if div != 0 {
						iterator.hour = mod
						iterator.day += div
						fixday = true
					}
				}
				if (len(r.byhour) == 0 || contains(r.byhour, iterator.hour)) &&
					(len(r.byminute) == 0 || contains(r.byminute, iterator.minute)) {
					break
				}
			}
			iterator.ii.fillTimeSet(&iterator.timeset, r.freq, iterator.hour, iterator.minute, iterator.second)
		} else if r.freq == SECONDLY {
			if filtered {
				// Jump to one iteration before next day
				iterator.second += (((86399 - (iterator.hour*3600 + iterator.minute*60 + iterator.second)) / r.interval) * r.interval)
			}
			for {
				iterator.second += r.interval
				div, mod := divmod(iterator.second, 60)
				if div != 0 {
					iterator.second = mod
					iterator.minute += div
					div, mod = divmod(iterator.minute, 60)
					if div != 0 {
						iterator.minute = mod
						iterator.hour += div
						div, mod = divmod(iterator.hour, 24)
						if div != 0 {
							iterator.hour = mod
							iterator.day += div
							fixday = true
						}
					}
				}
				if (len(r.byhour) == 0 || contains(r.byhour, iterator.hour)) &&
					(len(r.byminute) == 0 || contains(r.byminute, iterator.minute)) &&
					(len(r.bysecond) == 0 || contains(r.bysecond, iterator.second)) {
					break
				}
			}
			iterator.ii.fillTimeSet(&iterator.timeset, r.freq, iterator.hour, iterator.minute, iterator.second)
		}
		if fixday && iterator.day > 28 {
			daysinmonth := daysIn(iterator.month, iterator.year)
			if iterator.day > daysinmonth {
				for iterator.day > daysinmonth {
					iterator.day -= daysinmonth
					iterator.month++
					if iterator.month == 13 {
						iterator.month = 1
						iterator.year++
						if iterator.year > MAXYEAR {
							r.len = iterator.total
							iterator.finished = true
							return
						}
					}
					daysinmonth = daysIn(iterator.month, iterator.year)
				}
				iterator.ii.rebuild(iterator.year, iterator.month)
			}
		}
	}
}

func (iterator *rIterator) fillDaySetMonotonic(start, end int) {
	desiredLen := end - start

	if cap(iterator.dayset) < desiredLen {
		iterator.dayset = make([]optInt, 0, desiredLen)
	} else {
		iterator.dayset = iterator.dayset[:0]
	}

	for i := start; i < end; i++ {
		iterator.dayset = append(iterator.dayset, optInt{
			Int:     i,
			Defined: true,
		})
	}
}

// next returns next occurrence and true if it exists, else zero value and false
func (iterator *rIterator) next() (time.Time, bool) {
	iterator.generate()
	return iterator.remain.Pop()
}

type reusingRemainSlice struct {
	storage []time.Time
	backup  []time.Time
}

func (s reusingRemainSlice) Len() int {
	return len(s.storage)
}

func (s *reusingRemainSlice) Append(t time.Time) {
	s.storage = append(s.storage, t)
	s.backup = s.storage
}

func (s *reusingRemainSlice) Pop() (ret time.Time, ok bool) {
	if len(s.storage) == 0 {
		return time.Time{}, false
	}

	ret, s.storage = s.storage[0], s.storage[1:]

	if len(s.storage) == 0 {
		// flush storage
		s.storage = s.backup[:0]
	}

	return ret, true
}
