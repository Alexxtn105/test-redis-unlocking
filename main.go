package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// acquireLock Функция получения блокировки общего ресурса
func acquireLock(client *redis.Client, lockKey string, timeout time.Duration) bool {
	ctx := context.Background()

	// Пытаемся запросить блокировку командой SETNX (Set if Not eXists)
	lockAcquired, err := client.SetNX(ctx, lockKey, "1", timeout).Result()

	//если блокировка все еще есть, выходим с отрицательным результатом
	if err != nil {
		fmt.Println("Ошибка получения блокировки: ", err)
		return false
	}
	return lockAcquired
}

// releaseLock Функция разблокировки общего ресурса
func releaseLock(client *redis.Client, lockKey string) {
	ctx := context.Background()
	// удаляем ключ блокировки из redis
	client.Del(ctx, lockKey)
}

func main() {
	// Создаем клиент Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	//обязательно закрываем клиент redis
	defer client.Close()

	// объявляем имя блокировки и ее таймаут
	lockKey := "my_lock"
	lockTimeout := 10 * time.Second

	// получаем блокировку с указанными параметрами
	if acquireLock(client, lockKey, lockTimeout) {
		fmt.Println("Блокировка получена успешно!")

		// Симулируем длительный процесс
		time.Sleep(10 * time.Second)
		fmt.Println("Работа завершена")

		// Снимаем блокировку
		releaseLock(client, lockKey)
		fmt.Println("Блокировка снята")
	} else {
		fmt.Println("Ошибка получения блокировки. Ресурс уже заблокирован.")
	}

}
