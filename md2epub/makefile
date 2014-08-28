SOURCE = test
NAME = md2epub 

test: $(NAME)
	./$(NAME) $(SOURCE)
	rm -r out
	unzip -o $(SOURCE).epub -d out

$(NAME): build

build: 
	go build

.PHONY: test build