package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// Для элемента нужно знать его значение, предыдущий/следующий элементы.
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// Для списка нужно знать его длину, первый и последний элементы.
type list struct {
	length int
	first  *ListItem
	last   *ListItem
}

// Получить количество элементов в списке.
func (l *list) Len() int {
	return l.length
}

// Получить первый элемент из списка.
func (l *list) Front() *ListItem {
	return l.first
}

// Получить последний элемент из списка.
func (l *list) Back() *ListItem {
	return l.last
}

// Добавить значение в начало списка.
func (l *list) PushFront(v interface{}) *ListItem {
	// Готовим элемент. Значение - полученный. Следующий элемент будет тот который сейчас первый в списке.
	item := ListItem{Value: v, Next: l.first} // Если это первый вставляемый элемент - l.first будет nil
	// Элемент готов к вставке - следущий элемент (Next) вставили
	// Если это первый вставляемый элемент - тогда он же и является последним в списке
	if l.last == nil {
		l.last = &item
	}
	// Если первый элемент уже был - нам нужно обратиться к этому элементу
	// и добраться до параметра Prev, чтобы установить ему наш новый элемент
	if l.first != nil { // (это условно)
		l.first.Prev = &item
	}
	l.first = &item // сообщим списку информацию о новом первом элементе (это безусловно)

	l.length++
	return &item
}

// Добавить значение в конец списка.
func (l *list) PushBack(v interface{}) *ListItem {
	// Готовим элемент. Значение - полученный. Предыдущий элемент будет тот который сейчас последний в списке.
	item := ListItem{Value: v, Prev: l.last}
	// Если последний элемент в списке уже был - добираемся до него и в Next вставляем наш текущий элемент
	if l.last != nil {
		l.last.Next = &item
	}
	// Если у списка нет первого элемента - вставляемый элемент еще и является первым.
	if l.first == nil { // (это условно)
		l.first = &item
	}
	l.last = &item // сообщим списку информацию о новом последнем элементе (это безусловно)

	l.length++
	return &item
}

// Удаление элемента из списка.
func (l *list) Remove(i *ListItem) {
	// получаем информацию из элемента
	switch {
	case i.Prev != nil && i.Next != nil: // - значит наш элемент где-то посередине
		l.removeFromMiddle(i)
	case i.Prev == nil && i.Next != nil: // - значит наш элемент первый
		l.removeFromFront(i)
	case i.Prev != nil && i.Next == nil: // - значит наш элемент последний
		l.removeFromBack(i)
	case i.Prev == nil && i.Next == nil: // - значит наш элемент единственный в списке
		l.removeLast(i)
	}
}

// удаление из списка элемента, который где-то посередине
// проверяем что у предыдущего элемента, следующий записан наш
// и что у последующего элемента, предыдущий записан наш
// если так - то нужно предыдущему элементу установить следущий = текущему следующему
// и следущему элементу установить предыдущий = текущему предыдущему
func (l *list) removeFromMiddle(i *ListItem) {
	if i.Prev.Next == i && i.Next.Prev == i {
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
		l.length--
	}
}

// удаление из списка первого элемента
// смотрим какой элемент идет следующим - у него Prev должен быть равен нашему элементу
// если это так - нужно у этого следующего элемента стереть информацию о предыдущем(нашем) элементе
// и при этом сообщить списку о том что первым элементом теперь является тот, который ранее был вторым
// и уменьшаем length.
func (l *list) removeFromFront(i *ListItem) {
	if i.Next.Prev == i {
		i.Next.Prev, l.first = nil, i.Next
		l.length--
	}
}

// удаление из списка последнего элемента
// смотрим какой элемент был предыдущим - у него Next должен быть равен нашему элементу
// если это так - нужно у этого предыдущего элемента стереть информацию о следующем(нашем) элементе
// и при этом сообщить списку о том что последним элементом теперь является тот, который ранее был предпоследним
// и уменьшаем length.
func (l *list) removeFromBack(i *ListItem) {
	if i.Prev.Next == i {
		i.Prev.Next, l.last = nil, i.Prev
		l.length--
	}
}

// удаление из списка единственного элемента
// если элемент одновременно и первый и последний - удаляем информацию из переменных списка first/last
// и уменьшаем length до zero value (0).
func (l *list) removeLast(i *ListItem) {
	if l.first == i && l.last == i {
		l.first, l.last = nil, nil
		l.length--
	}
}

// переместить элемент в начало.
// Полученный элемент удаляем иеющимся методом
// Далее идет код аналогичный вставке элемента в началою.
func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	item := ListItem{Next: l.first, Value: i.Value}
	if l.last == nil {
		l.last = &item
	}
	if l.first != nil {
		l.first.Prev = &item
	}
	l.first = &item
	l.length++
}

func NewList() *list {
	return new(list)
}
