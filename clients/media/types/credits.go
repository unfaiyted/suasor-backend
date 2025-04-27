package types

type Credits []Person

func (c Credits) GetCast() []Person {
	var cast []Person
	for _, person := range c {
		if person.IsCast {
			cast = append(cast, person)
		}
	}
	return cast
}

func (c Credits) GetCrew() []Person {
	var crew []Person
	for _, person := range c {
		if person.IsCrew {
			crew = append(crew, person)
		}
	}
	return crew
}

func (c Credits) GetGuests() []Person {
	var guests []Person
	for _, person := range c {
		if person.IsGuest {
			guests = append(guests, person)
		}
	}
	return guests
}

func (c Credits) GetCreators() []Person {
	var creators []Person
	for _, person := range c {
		if person.IsCreator {
			creators = append(creators, person)
		}
	}
	return creators
}
