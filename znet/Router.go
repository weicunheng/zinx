package znet

/*
实现Router时，先嵌入BaseRouter基类，然后再根据自身业务的需要，重写需要的方法即可
*/
type BaseRouter struct {
}

func (r *BaseRouter) PreHandler() {

}
func (r *BaseRouter) Handler(){

}
func (r *BaseRouter) PostHandler(){

}