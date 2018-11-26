package session

//
//type RedisSessionService struct {
//	Conn redis.Conn
//}
//
//
//func (r *RedisSessionService) InitService() (redis.Conn, error) {
//	port := "6379"
//	var err error
//	r.Conn, err = redis.Dial("tcp", "localhost:" + port)
//	if err != nil {
//		fmt.Println(err.Error())
//		fmt.Println("Error in start Redis")
//		return nil, err
//	}
//	//defer r.Conn.Close()
//	fmt.Println("Redis start in port: ", port)
//	return r.Conn, nil
//}
//
//func (r *RedisSessionService) Create(key string, value string) (error){
//	_, err := r.Conn.Do("SET", key, value)
//
//	if err != nil {
//		return err
//	}
//
//	//defer r.Conn.Close()
//
//	return nil
//}
//
//func (r *RedisSessionService) Get(key string) (string, error){
//	data, err := r.Conn.Do("GET", key)
//	item, err := redis.String(data, err)
//
//	if err != nil {
//		if err == redis.ErrNil {
//			return "", redis.ErrNil
//		} else {
//			return "", err
//		}
//	}
//
//	//defer r.Conn.Close()
//
//	return item, nil
//
//}
//
//func (r *RedisSessionService) Delete(key string) (error){
//	_, err := r.Conn.Do("DEL", key)
//	//item, err := redis.String(data, err)
//
//	if err != nil {
//		return err
//	}
//
//	//defer r.Conn.Close()
//
//	return nil
//
//}

//func (r *RedisSessionService) Delete(user *models.User) (error) {
//
//	_, err := r.Conn.Do("DEL", key)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

