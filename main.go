package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type students_complaint struct {
	Uid string `json:"uid"`
	Roll_No string `json:"roll_no"`
	Complaint_Text string `json:"complaint_text"`
	Complaint_Text_Title string `json:"complaint_text_title"`
	Hostel_Code string `json:"hostel_code"`
	Query_Resolved string `json:"query_resolved"`
}

type admin_database struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Password string `json:"password"`
	Hostel_Code string `json:"hostel_code"`
}

type students_database struct {
    Roll_No string  `json:"roll_no"`
    Name  string  `json:"name"`
    Password string `json:"password"`
    Hostel_Code string  `json:"hostel_code"`
}

var Students_Data []students_database
var Admin_Data []admin_database
var Complaint_Data []students_complaint

func GatherAdminData(c *gin.Context) {
	username := c.Param("username")

    for _, a := range Admin_Data {
        if a.Username == username {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
	}
}

func NewComplaint(c *gin.Context) {
	var complaint_data students_complaint

	if c.BindJSON(&complaint_data) != nil {
        return
    }
    
    AddComplaintToDatabase(complaint_data)
    Complaint_Data = append(Complaint_Data, complaint_data)
}

func GatherHostelStudentsData(c *gin.Context) {
	hostel_code := c.Param("hostel_code")

	for _, a := range Students_Data {
        if a.Hostel_Code == hostel_code {
            c.IndentedJSON(http.StatusOK, a)
            break
        }
	}
}

func GatherStudentsComplaints(c *gin.Context) {
    hostel_code := c.Param("hostel_code")

	for _, a := range Complaint_Data {
        if a.Hostel_Code == hostel_code {
            c.IndentedJSON(http.StatusOK, a)
        }
	}
}

func GatherUserData(c *gin.Context) {
	roll_no := c.Param("roll_no")

    for _, a := range Students_Data {
        if a.Roll_No == roll_no {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
	}
}

func GatherUserComplaints(c *gin.Context) {
	roll_no := c.Param("roll_no")

	for _, a := range Complaint_Data {
        if a.Roll_No == roll_no {
            c.IndentedJSON(http.StatusOK, a)
        }
	}
}

func ResolveUserComplaint(c *gin.Context) {
	uid := c.Param("uid")
	var complaint_data students_complaint

	if c.BindJSON(&complaint_data) != nil {
        return
    }

    for i, a := range Complaint_Data {
        if a.Roll_No == complaint_data.Roll_No && a.Uid == uid{
            Complaint_Data[i].Query_Resolved = complaint_data.Query_Resolved
            UserComplaintResolve(complaint_data.Query_Resolved, uid)
            return
        }
	}
}

func main() {
	OpenDatabase()
	GatherDataFromDatabase()

	r := engine()
	r.Use(gin.Logger())
	if err := engine().Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	r := gin.New()

	// Login and logout routes
	r.POST("/login", LoginFunc)
	r.POST("/admin/login", AdminLoginFunc)
	r.POST("/logout", LogoutFunc)
	r.POST("/admin/logout", AdminLogoutFunc)
	r.POST("/signin", SigninUser)

	// Private group, require authentication to access
	private := r.Group("/private")
	private.Use(AuthRequired)
	{
		private.POST("/user/:roll_no", GatherUserData)
		private.POST("/complaint", NewComplaint)
		private.POST("/complaints:roll_no", GatherUserComplaints)
		private.POST("/complaint/:uid/resolve", ResolveUserComplaint)
	}

	admin := r.Group("/admin")
	admin.Use(AdminAuthRequired)
	{
		admin.POST("/admin/:username", GatherAdminData)
		admin.POST("/users:hostel_code", GatherHostelStudentsData)
		admin.POST("/complaints:hostel_code", GatherStudentsComplaints)
		admin.POST("/complaint/:uid/resolve", ResolveUserComplaint)
	}
	return r
}
