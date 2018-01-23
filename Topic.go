package main
import(
	"strings"
)

var UsersTopocs map[string]map[string]string
var Topics map[string]*topic

func ListenTopic(userid, topicName string) {
	topicName = strings.ToUpper(topicName)
	thisTopic, isTrue := Topics[topicName]
	if !isTrue {
		thisTopic = &topic{Name: topicName, Users: make(map[string]string)}
		Topics[topicName] = thisTopic
	}
	thisTopic.Users[userid] = userid

	_ ,isTrue=UsersTopocs[userid]
	if !isTrue{
		UsersTopocs[userid] = make(map[string]string)
	}
	UsersTopocs[userid][topicName]=topicName
	
}
func UnsubscribeUserAll(userid string) {
	for _,v := range UsersTopocs[userid]{
		Unsubscribe(userid,v)
	}
}
func Unsubscribe(userid, topicName string) {
	topicName = strings.ToUpper(topicName)
	delete(UsersTopocs[userid],topicName)
	if len(UsersTopocs[userid]) <1{
		delete(UsersTopocs,userid)
	}
	delete(Topics[topicName].Users,userid)
	if len(Topics[topicName].Users) <1{
		delete(Topics,topicName)
	}

}


